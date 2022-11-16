package utils

import (
	grpcClient "bro-n-bro-osmosis/client/grpc"
	"bro-n-bro-osmosis/types"
	"context"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// StoreValidatorFromMsgCreateValidator handles properly a MsgCreateValidator instance by
// saving into the database all the data associated to such validator
func StoreValidatorFromMsgCreateValidator(height int64, msg *stakingtypes.MsgCreateValidator, cdc codec.Codec) error {
	var pubKey cryptotypes.PubKey
	err := cdc.UnpackAny(msg.Pubkey, &pubKey)
	if err != nil {
		return err
	}

	operatorAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return err
	}

	stakingValidator, err := stakingtypes.NewValidator(operatorAddr, pubKey, msg.Description)
	if err != nil {
		return err
	}

	validator, err := ConvertValidator(cdc, stakingValidator, height)
	if err != nil {
		return err
	}

	desc, err := ConvertValidatorDescription(msg.ValidatorAddress, msg.Description, height)
	if err != nil {
		return err
	}

	_, _ = validator, desc
	// TODO:!!!!!
	// Save the validator
	//err = db.SaveValidatorsData([]types.Validator{validator})
	//if err != nil {
	//	return err
	//}
	//
	// Save the description
	//err = db.SaveValidatorDescription(desc)
	//if err != nil {
	//	return err
	//}
	//
	// Save the first self-delegation
	//err = db.SaveDelegations([]types.Delegation{
	//	types.NewDelegation(
	//		msg.DelegatorAddress,
	//		msg.ValidatorAddress,
	//		msg.Value,
	//		height,
	//	),
	//})
	//if err != nil {
	//	return err
	//}

	// Save the commission
	//err = db.SaveValidatorCommission(types.NewValidatorCommission(
	//	msg.ValidatorAddress,
	//	&msg.Commission.Rate,
	//	&msg.MinSelfDelegation,
	//	height,
	//))
	return err
}

// StoreDelegationFromMessage handles a MsgDelegate and saves the delegation inside the database
func StoreDelegationFromMessage(height int64, msg *stakingtypes.MsgDelegate, stakingClient stakingtypes.QueryClient) error {
	header := grpcClient.GetHeightRequestHeader(height)
	res, err := stakingClient.Delegation(
		context.Background(),
		&stakingtypes.QueryDelegationRequest{
			DelegatorAddr: msg.DelegatorAddress,
			ValidatorAddr: msg.ValidatorAddress,
		},
		header,
	)
	if err != nil {
		return err
	}

	delegation := ConvertDelegationResponse(height, *res.DelegationResponse)
	//return db.SaveDelegations([]types.Delegation{delegation})
	// TODO:
	_ = delegation
	return nil
}

// StoreRedelegationFromMessage handles a MsgBeginRedelegate by saving the redelegation inside the database,
// and returns the new redelegation instance
func StoreRedelegationFromMessage(tx *types.Tx, index int, msg *stakingtypes.MsgBeginRedelegate) (*types.Redelegation, error) {
	event, err := tx.FindEventByType(index, stakingtypes.EventTypeRedelegate)
	if err != nil {
		return nil, err
	}

	completionTimeStr, err := tx.FindAttributeByKey(event, stakingtypes.AttributeKeyCompletionTime)
	if err != nil {
		return nil, err
	}

	completionTime, err := time.Parse(time.RFC3339, completionTimeStr)
	if err != nil {
		return nil, err
	}

	redelegation := types.NewRedelegation(
		msg.DelegatorAddress,
		msg.ValidatorSrcAddress,
		msg.ValidatorDstAddress,
		msg.Amount,
		completionTime,
		tx.Height,
	)

	// TODO:
	//err = db.SaveRedelegations([]types.Redelegation{redelegation})

	return &redelegation, err
}

// StoreUnbondingDelegationFromMessage handles a MsgUndelegate storing the new unbonding delegation inside the database,
// and returns the new unbonding delegation instance
func StoreUnbondingDelegationFromMessage(tx *types.Tx, index int, msg *stakingtypes.MsgUndelegate) (*types.UnbondingDelegation, error) {
	event, err := tx.FindEventByType(index, stakingtypes.EventTypeUnbond)
	if err != nil {
		return nil, err
	}

	completionTimeStr, err := tx.FindAttributeByKey(event, stakingtypes.AttributeKeyCompletionTime)
	if err != nil {
		return nil, err
	}

	completionTime, err := time.Parse(time.RFC3339, completionTimeStr)
	if err != nil {
		return nil, err
	}

	delegation := types.NewUnbondingDelegation(
		msg.DelegatorAddress,
		msg.ValidatorAddress,
		msg.Amount,
		completionTime,
		tx.Height,
	)

	// TODO:
	//err =  db.SaveUnbondingDelegations([]types.UnbondingDelegation{delegation})

	return &delegation, err
}
