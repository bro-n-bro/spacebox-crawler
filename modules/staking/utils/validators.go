package utils

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"bro-n-bro-osmosis/adapter/broker/model"
	grpcClient "bro-n-bro-osmosis/client/grpc"
	"bro-n-bro-osmosis/internal/rep"
	"bro-n-bro-osmosis/modules/staking/keybase"
	tb "bro-n-bro-osmosis/pkg/mapper/to_broker"
	"bro-n-bro-osmosis/types"
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

// ---------------------------------------------------------------------------------------------------------------------

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

	var operator sdk.ValAddress
	operator = validator.GetOperator() // FIXME: here was a panic: invalid Bech32 prefix; expected cosmosvaloper, got bostromvaloper
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

// ConvertValidatorDescription returns a new types.ValidatorDescription object by fetching the avatar URL
// using the Keybase APIs
func ConvertValidatorDescription(
	opAddr string, description stakingtypes.Description, height int64,
) (types.ValidatorDescription, error) {
	avatarURL, err := keybase.GetAvatarURL(description.Identity)
	if err != nil {
		return types.ValidatorDescription{}, err
	}

	return types.NewValidatorDescription(opAddr, description, avatarURL, height), nil
}

// --------------------------------------------------------------------------------------------------------------------

// GetValidators returns the validators list at the given height
func GetValidators(height int64, stakingClient stakingtypes.QueryClient, cdc codec.Codec,
) ([]stakingtypes.Validator, []types.StakingValidator, error) {

	return GetValidatorsWithStatus(height, "", stakingClient, cdc)
}

// GetValidatorsWithStatus returns the list of all the validators having the given status at the given height
func GetValidatorsWithStatus(
	height int64, status string, stakingClient stakingtypes.QueryClient, cdc codec.Codec,
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

// UpdateValidators updates the list of validators that are present at the given height
func UpdateValidators(ctx context.Context, height int64, client stakingtypes.QueryClient, cdc codec.Codec,
	broker rep.Broker, mapper tb.ToBroker) ([]stakingtypes.Validator, error) {

	vals, validators, err := GetValidators(height, client, cdc)
	if err != nil {
		return nil, err
	}

	// TODO: save to mongo?
	// TODO: test it
	if err = publishValidatorsData(ctx, validators, broker, mapper); err != nil {
		return nil, err
	}

	//err = db.SaveValidatorsData(validators)
	//if err != nil {
	//	return nil, err
	//}

	// TODO:
	_ = validators

	return vals, err
}

// --------------------------------------------------------------------------------------------------------------------

func GetValidatorsStatuses(height int64, validators []stakingtypes.Validator, cdc codec.Codec) ([]types.ValidatorStatus, error) {
	statuses := make([]types.ValidatorStatus, len(validators))
	for index, validator := range validators {
		consAddr, err := GetValidatorConsAddr(cdc, validator)
		if err != nil {
			return nil, fmt.Errorf("error while getting validator consensus address: %s", err)
		}

		consPubKey, err := GetValidatorConsPubKey(cdc, validator)
		if err != nil {
			return nil, fmt.Errorf("error while getting validator consensus public key: %s", err)
		}

		statuses[index] = types.NewValidatorStatus(
			consAddr.String(),
			consPubKey.String(),
			int(validator.GetStatus()),
			validator.IsJailed(),
			height,
		)
	}

	return statuses, nil
}

func GetValidatorsVotingPowers(height int64, vals *tmctypes.ResultValidators) []types.ValidatorVotingPower {
	votingPowers := make([]types.ValidatorVotingPower, len(vals.Validators))
	for index, validator := range vals.Validators {
		consAddr := sdk.ConsAddress(validator.Address).String()
		// FIXME: how to check it?
		//if found, _ := db.HasValidator(consAddr); !found {
		//	continue
		//}

		votingPowers[index] = types.NewValidatorVotingPower(consAddr, validator.VotingPower, height)
	}
	return votingPowers
}

func publishValidatorsData(ctx context.Context, sVals []types.StakingValidator, broker rep.Broker,
	mapper tb.ToBroker) error {

	vals := make(types.Validators, len(sVals))
	accounts := make([]types.Account, len(sVals))
	valsInfo := make([]model.ValidatorInfo, len(sVals))
	for i, val := range sVals {
		vals[i] = types.NewValidator(val.GetConsAddr(), val.GetConsPubKey())
		accounts[i] = types.NewAccount(val.GetSelfDelegateAddress(), val.GetHeight())
		valsInfo[i] = mapper.MapValidatorInfo(val)
	}

	err := broker.PublishValidators(ctx, mapper.MapValidators(vals))
	if err != nil {
		return err
	}
	err = broker.PublishAccounts(ctx, mapper.MapAccounts(accounts))
	if err != nil {
		return err
	}
	err = broker.PublishValidatorsInfo(ctx, valsInfo)
	if err != nil {
		return err
	}

	return nil
}
