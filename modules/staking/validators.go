package staking

import (
	"context"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/bro-n-bro/spacebox-crawler/pkg/keybase"
	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

const (
	defaultLimit = 150
)

type avatarURLCache struct {
	Identity  string
	AvatarURL string
}

// HandleValidators handles validators for each block height.
func (m *Module) HandleValidators(ctx context.Context, tmValidators *tmctypes.ResultValidators) error {
	vals, validators, err := GetValidators(ctx, tmValidators.BlockHeight, m.client.StakingQueryClient, m.cdc)
	if err != nil {
		return err
	}

	if err = m.PublishValidatorsData(ctx, validators); err != nil {
		return err
	}

	if err = m.publishValidatorDescriptions(ctx, vals, tmValidators.BlockHeight, m.parseAvatarURL); err != nil {
		return err
	}

	for _, val := range vals {
		consAddr, err := getValidatorConsAddr(m.cdc, val)
		if err != nil {
			return fmt.Errorf("error while getting validator consensus address: %w", err)
		}

		if err = m.broker.PublishValidatorStatus(ctx, model.ValidatorStatus{
			Height:           tmValidators.BlockHeight,
			ValidatorAddress: consAddr.String(),
			Status:           int64(val.GetStatus()),
			Jailed:           val.IsJailed(),
		}); err != nil {
			return err
		}

		if err := m.broker.PublishValidatorCommission(ctx, model.ValidatorCommission{
			Height:          tmValidators.BlockHeight,
			OperatorAddress: val.OperatorAddress,
			Commission:      val.Commission.Rate.MustFloat64(),
			MaxChangeRate:   val.Commission.MaxChangeRate.MustFloat64(),
			MaxRate:         val.Commission.MaxRate.MustFloat64(),
		}); err != nil {
			return err
		}
	}

	// FIXME: does it needed?
	// Update the voting powers
	// go updateValidatorVotingPower(block.Height, vals)

	return nil
}

// PublishValidatorsData produces a message about validator, account and validator info to the broker.
func (m *Module) PublishValidatorsData(ctx context.Context, sVals []types.StakingValidator) error {
	prefix := sdk.GetConfig().GetBech32AccountAddrPrefix()

	for _, val := range sVals {
		// TODO: test it
		if err := m.broker.PublishValidator(ctx, model.Validator{
			ConsensusAddress: val.GetConsAddr(),
			ConsensusPubkey:  val.GetConsPubKey(),
			OperatorAddress:  val.GetOperator(),
			Height:           val.GetHeight(),
		}); err != nil {
			return err
		}

		if strings.HasPrefix(val.GetSelfDelegateAddress(), prefix) {
			// TODO: test it
			if err := m.broker.PublishAccount(ctx, model.Account{
				Address: val.GetSelfDelegateAddress(),
				Height:  val.GetHeight(),
			}); err != nil {
				return err
			}
		}

		var minSelfDelegation int64
		if val.GetMinSelfDelegation() != nil {
			minSelfDelegation = val.GetMinSelfDelegation().Int64()
		}

		// TODO: test it
		if err := m.broker.PublishValidatorInfo(ctx, model.ValidatorInfo{
			ConsensusAddress:    val.GetConsAddr(),
			OperatorAddress:     val.GetOperator(),
			SelfDelegateAddress: val.GetSelfDelegateAddress(),
			MinSelfDelegation:   minSelfDelegation,
			Height:              val.GetHeight(),
		}); err != nil {
			return err
		}
	}

	return nil
}

// asyncPublishValidatorDescriptions process validator descriptions and publish them to the broker.
func (m *Module) publishValidatorDescriptions(ctx context.Context, vals stakingtypes.Validators, height int64, parseAvatarURL bool) error {
	for _, val := range vals {
		if parseAvatarURL {
			go func(val stakingtypes.Validator) {
				ctx = context.Background()
				avatarURL := m.getAvatarURL(val.OperatorAddress, val.Description.Identity, height)
				_ = m.publishValidatorDescription(ctx, val, avatarURL, height)
			}(val)
		} else {
			if err := m.publishValidatorDescription(ctx, val, "", height); err != nil {
				return err
			}
		}
	}

	return nil
}

// publishValidatorDescription process validator description and publish it to the broker.
// It also gets avatar url from the keybase.
// Contains the cache for validator identity to skip the keybase API call as it might be stopped due to rate limits.
func (m *Module) publishValidatorDescription(ctx context.Context, val stakingtypes.Validator, avatarURL string, height int64) error {
	if err := m.broker.PublishValidatorDescription(ctx, model.ValidatorDescription{
		OperatorAddress: val.OperatorAddress,
		Moniker:         val.Description.Moniker,
		Identity:        val.Description.Identity,
		Website:         val.Description.Website,
		SecurityContact: val.Description.SecurityContact,
		Details:         val.Description.Details,
		AvatarURL:       avatarURL,
		Height:          height,
	}); err != nil {
		m.log.Error().Err(err).
			Str("operator_address", val.OperatorAddress).
			Str("identity", val.Description.Identity).
			Int64("height", height).
			Msg("failed to publish validator description")
		return err
	}

	return nil
}

func (m *Module) getAvatarURL(operatorAddress, identity string, height int64) string {
	var (
		avatarURL string
		cacheItem avatarURLCache
		err       error
		ctx       = context.Background()
	)

	cacheVal, ok := m.avatarURLCache.Load(operatorAddress)
	if ok {
		cacheItem, ok = cacheVal.(avatarURLCache)
	}

	// not exists or value is not equal to the current one
	if !ok || cacheItem.Identity != identity {
		// get avatar url from the keybase API
		avatarURL, err = keybase.GetAvatarURL(ctx, identity)
		if err != nil {
			m.log.Error().
				Err(err).
				Str("operator_address", operatorAddress).
				Str("identity", identity).
				Int64("height", height).
				Msg("failed to get avatar url")
		} else {
			// update the cache
			m.avatarURLCache.Store(operatorAddress, avatarURLCache{
				Identity:  identity,
				AvatarURL: avatarURL,
			})
		}
	} else { // can get from cache
		avatarURL = cacheItem.AvatarURL
	}

	return avatarURL
}
