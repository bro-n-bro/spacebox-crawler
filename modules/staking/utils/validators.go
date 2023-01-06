package utils

import (
	"context"

	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	grpcClient "github.com/hexy-dev/spacebox-crawler/client/grpc"
	"github.com/hexy-dev/spacebox-crawler/types"
	"github.com/hexy-dev/spacebox/broker/model"
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

	operator := validator.GetOperator() // FIXME: here was a panic: invalid Bech32 prefix; expected cosmosvaloper, got bostromvaloper
	return types.NewStakingValidator(
		consAddr.String(),
		validator.OperatorAddress,
		consPubKey.String(),
		sdk.AccAddress(operator).String(),
		&validator.Commission.MaxChangeRate,
		&validator.Commission.MaxRate,
		height,
	), nil
}

// GetValidators returns the validators list at the given height
func GetValidators(height int64, stakingClient stakingtypes.QueryClient, cdc codec.Codec,
) ([]stakingtypes.Validator, []types.StakingValidator, error) {

	return GetValidatorsWithStatus(height, "", stakingClient, cdc)
}

// GetValidatorsWithStatus returns the list of all the validators having the given status at the given height.
func GetValidatorsWithStatus(height int64, status string, stakingClient stakingtypes.QueryClient, cdc codec.Codec,
) ([]stakingtypes.Validator, []types.StakingValidator, error) {

	header := grpcClient.GetHeightRequestHeader(height)

	var validators []stakingtypes.Validator
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := stakingClient.Validators(
			context.Background(),
			&stakingtypes.QueryValidatorsRequest{
				Status: status,
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 validators at time
				},
			},
			header,
		)
		if err != nil {
			return nil, nil, err
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		validators = append(validators, res.Validators...)
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
		PublishAccounts(ctx context.Context, accounts []model.Account) error // FIXME: auth module
		PublishValidators(ctx context.Context, vals []model.Validator) error
		PublishValidatorsInfo(ctx context.Context, infos []model.ValidatorInfo) error
	},
) ([]stakingtypes.Validator, error) {

	vals, validators, err := GetValidators(height, client, cdc)
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
		PublishAccounts(ctx context.Context, accounts []model.Account) error // FIXME: auth module
		PublishValidators(ctx context.Context, vals []model.Validator) error
		PublishValidatorsInfo(ctx context.Context, infos []model.ValidatorInfo) error
	},
) error {

	for _, val := range sVals {
		// TODO: test it
		err := broker.PublishValidators(ctx, []model.Validator{model.NewValidator(val.GetConsAddr(), val.GetConsPubKey())})
		if err != nil {
			return err
		}

		// TODO: test it
		err = broker.PublishAccounts(ctx, []model.Account{model.NewAccount(val.GetSelfDelegateAddress(), val.GetHeight())})
		if err != nil {
			return err
		}

		var minSelfDelegation int64
		if val.GetMinSelfDelegation() != nil {
			minSelfDelegation = val.GetMinSelfDelegation().Int64()
		}

		vi := model.NewValidatorInfo(
			val.GetHeight(),
			minSelfDelegation,
			val.GetConsAddr(),
			val.GetOperator(),
			val.GetSelfDelegateAddress(),
		)

		// TODO: test it
		err = broker.PublishValidatorsInfo(ctx, []model.ValidatorInfo{vi})
		if err != nil {
			return err
		}
	}

	return nil
}
