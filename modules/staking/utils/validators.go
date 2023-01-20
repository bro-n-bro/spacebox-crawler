package utils

import (
	"context"

	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

const (
	defaultLimit = 150
)

// GetValidatorConsPubKey returns the consensus public key of the given validator
func GetValidatorConsPubKey(cdc codec.Codec, validator stakingtypes.Validator) (cryptotypes.PubKey, error) {
	var pubKey cryptotypes.PubKey
	err := cdc.UnpackAny(validator.ConsensusPubkey, &pubKey)

	return pubKey, err
}

// GetValidatorConsAddr returns the consensus address of the given validator
func GetValidatorConsAddr(cdc codec.Codec, validator stakingtypes.Validator) (sdk.ConsAddress, error) {
	pubKey, err := GetValidatorConsPubKey(cdc, validator)
	if err != nil {
		return nil, err
	}

	return sdk.ConsAddress(pubKey.Address()), err
}

// ConvertValidator converts the given staking validator into a BDJuno validator
func ConvertValidator(cdc codec.Codec, validator stakingtypes.Validator, height int64) (types.StakingValidator, error) {
	consAddr, err := GetValidatorConsAddr(cdc, validator)
	if err != nil {
		return nil, err
	}

	consPubKey, err := GetValidatorConsPubKey(cdc, validator)
	if err != nil {
		return nil, err
	}

	return types.NewStakingValidator(
		consAddr.String(),
		validator.OperatorAddress,
		consPubKey.String(),
		sdk.AccAddress(validator.GetOperator()).String(),
		&validator.Commission.MaxChangeRate,
		&validator.Commission.MaxRate,
		height,
	), nil
}

// GetValidators returns the validators list at the given height
func GetValidators(ctx context.Context, height int64, stakingClient stakingtypes.QueryClient,
	cdc codec.Codec) ([]stakingtypes.Validator, []types.StakingValidator, error) {

	return GetValidatorsWithStatus(ctx, height, "", stakingClient, cdc)
}

// GetValidatorsWithStatus returns the list of all the validators having the given status at the given height.
func GetValidatorsWithStatus(ctx context.Context, height int64, status string, stakingClient stakingtypes.QueryClient,
	cdc codec.Codec) ([]stakingtypes.Validator, []types.StakingValidator, error) {

	header := grpcClient.GetHeightRequestHeader(height)

	var (
		validators []stakingtypes.Validator
		nextKey    []byte
	)

	for {
		respPb, err := stakingClient.Validators(
			ctx,
			&stakingtypes.QueryValidatorsRequest{
				Status: status,
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: defaultLimit,
				},
			},
			header,
		)
		if err != nil {
			return nil, nil, err
		}

		if len(nextKey) == 0 { // first iteration
			validators = make([]stakingtypes.Validator, 0, respPb.Pagination.Total)
		}

		nextKey = respPb.Pagination.NextKey
		validators = append(validators, respPb.Validators...)

		if len(respPb.Pagination.NextKey) == 0 {
			break
		}
	}

	var vals = make([]types.StakingValidator, len(validators))

	for index, val := range validators {
		validator, err := ConvertValidator(cdc, val, height)
		if err != nil {
			return nil, nil, err
		}

		vals[index] = validator
	}

	return validators, vals, nil
}

// UpdateValidators updates the list of validators that are present at the given height and produces it to the broker.
func UpdateValidators(
	ctx context.Context,
	height int64,
	client stakingtypes.QueryClient,
	cdc codec.Codec,
	broker interface {
		PublishAccount(ctx context.Context, accounts model.Account) error // FIXME: auth module
		PublishValidator(ctx context.Context, val model.Validator) error
		PublishValidatorInfo(ctx context.Context, info model.ValidatorInfo) error
	},
) ([]stakingtypes.Validator, error) {

	vals, validators, err := GetValidators(ctx, height, client, cdc)
	if err != nil {
		return nil, err
	}

	// TODO: save to mongo?
	// TODO: test it
	if err = PublishValidatorsData(ctx, validators, broker); err != nil {
		return nil, err
	}

	return vals, err
}

// PublishValidatorsData produces a message about validator, account and validator info for the broker.
func PublishValidatorsData(
	ctx context.Context,
	sVals []types.StakingValidator,
	broker interface {
		PublishAccount(ctx context.Context, accounts model.Account) error // FIXME: auth module
		PublishValidator(ctx context.Context, val model.Validator) error
		PublishValidatorInfo(ctx context.Context, info model.ValidatorInfo) error
	},
) error {

	for _, val := range sVals {
		// TODO: test it
		if err := broker.PublishValidator(ctx, model.Validator{
			ConsensusAddress: val.GetConsAddr(),
			ConsensusPubkey:  val.GetConsPubKey(),
		}); err != nil {
			return err
		}

		// TODO: test it
		if err := broker.PublishAccount(ctx, model.Account{
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
		if err := broker.PublishValidatorInfo(ctx, model.ValidatorInfo{
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
