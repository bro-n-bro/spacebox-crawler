package staking

import (
	"context"
	"fmt"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	"golang.org/x/sync/errgroup"

	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

const (
	defaultLimit = 150
)

// HandleValidators handles validators for each block height.
func (m *Module) HandleValidators(ctx context.Context, vals *tmctypes.ResultValidators) error {
	// Update the validators
	validators, err := m.UpdateValidators(ctx, vals.BlockHeight)
	if err != nil {
		return err
	}
	g, ctx2 := errgroup.WithContext(ctx)

	g.Go(func() error {
		// 	Update the validators statuses
		return m.updateValidatorsStatus(ctx2, vals.BlockHeight, validators)
	})

	// FIXME: does it needed?
	// Update the voting powers
	// go updateValidatorVotingPower(block.Height, vals)

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

// updateValidatorsStatus updates all validators statuses
func (m *Module) updateValidatorsStatus(ctx context.Context, height int64, vals []stakingtypes.Validator) error {
	for _, validator := range vals {
		consAddr, err := getValidatorConsAddr(m.cdc, validator)
		if err != nil {
			return fmt.Errorf("error while getting validator consensus address: %w", err)
		}

		// consPubKey, err := stakingutils.getValidatorConsPubKey(m.cdc, validator)
		// if err != nil {
		// 	return fmt.Errorf("error while getting validator consensus public key: %w", err)
		// }

		// TODO: save to mongo?
		// TODO: test it
		// if err = m.broker.PublishValidator(ctx, model.Validator{
		// 	ConsensusAddress: consAddr.String(),
		// 	ConsensusPubkey:  consPubKey.String(),
		// }); err != nil {
		// 	return err
		// }

		// TODO: test it
		if err = m.broker.PublishValidatorStatus(ctx, model.ValidatorStatus{
			Height:           height,
			ValidatorAddress: consAddr.String(),
			Status:           int64(validator.GetStatus()),
			Jailed:           validator.IsJailed(),
		}); err != nil {
			return err
		}
	}

	return nil
}

// UpdateValidators updates the list of validators that are present at the given height and produces it to the broker.
func (m *Module) UpdateValidators(ctx context.Context, height int64) ([]stakingtypes.Validator, error) {
	vals, validators, err := GetValidators(ctx, height, m.client.StakingQueryClient, m.cdc)
	if err != nil {
		return nil, err
	}

	// TODO: save to mongo?
	// TODO: test it
	if err = m.PublishValidatorsData(ctx, validators); err != nil {
		return nil, err
	}

	return vals, err
}

// PublishValidatorsData produces a message about validator, account and validator info to the broker.
func (m *Module) PublishValidatorsData(ctx context.Context, sVals []types.StakingValidator) error {
	for _, val := range sVals {
		// TODO: test it
		if err := m.broker.PublishValidator(ctx, model.Validator{
			ConsensusAddress: val.GetConsAddr(),
			ConsensusPubkey:  val.GetConsPubKey(),
			OperatorAddress:  val.GetOperator(),
		}); err != nil {
			return err
		}

		// TODO: test it
		if err := m.broker.PublishAccount(ctx, model.Account{
			Address: val.GetSelfDelegateAddress(),
			Height:  val.GetHeight(),
		}); err != nil {
			return err
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
