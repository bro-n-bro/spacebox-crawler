package staking

import (
	"context"

	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	"github.com/bro-n-bro/spacebox-crawler/types"
)

// GetValidators returns the validators list at the given height.
func GetValidators(ctx context.Context, height int64, stakingClient stakingtypes.QueryClient,
	cdc codec.Codec) ([]stakingtypes.Validator, []types.StakingValidator, error) {

	return GetStakingValidators(ctx, height, "", stakingClient, cdc)
}

// GetStakingValidators returns the list of all the validators having the given status at the given height.
func GetStakingValidators(ctx context.Context, height int64, status string, stakingClient stakingtypes.QueryClient,
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

	// mapping by pubkey
	// tm validator address == consensus address
	// staking validator pubkey == tm validator pubkey

	var vals = make([]types.StakingValidator, len(validators))

	for index, val := range validators {
		validator, err := convertValidator(cdc, val, height)
		if err != nil {
			return nil, nil, err
		}

		vals[index] = validator
	}

	return validators, vals, nil
}

// getValidatorConsPubKey returns the consensus public key of the given validator.
func getValidatorConsPubKey(cdc codec.Codec, validator stakingtypes.Validator) (cryptotypes.PubKey, error) {
	var pubKey cryptotypes.PubKey
	err := cdc.UnpackAny(validator.ConsensusPubkey, &pubKey)

	return pubKey, err
}

// getValidatorConsAddr returns the consensus address of the given validator.
func getValidatorConsAddr(cdc codec.Codec, validator stakingtypes.Validator) (sdk.ConsAddress, error) {
	pubKey, err := getValidatorConsPubKey(cdc, validator)
	if err != nil {
		return nil, err
	}

	return sdk.ConsAddress(pubKey.Address()), err
}

// convertValidator converts the given staking validator in types.StakingValidator instance.
func convertValidator(cdc codec.Codec, validator stakingtypes.Validator, height int64) (types.StakingValidator, error) {
	consAddr, err := getValidatorConsAddr(cdc, validator)
	if err != nil {
		return nil, err
	}

	consPubKey, err := getValidatorConsPubKey(cdc, validator)
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
		validator.Description,
		height,
	), nil
}
