package utils

import (
	"context"
	"time"

	grpcClient "github.com/hexy-dev/spacebox-crawler/client/grpc"
	"github.com/hexy-dev/spacebox-crawler/internal/rep"
	tb "github.com/hexy-dev/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/hexy-dev/spacebox-crawler/types"

	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// StoreValidatorFromMsgCreateValidator handles properly a MsgCreateValidator instance by
// saving into the database all the data associated to such validator
func StoreValidatorFromMsgCreateValidator(ctx context.Context, height int64, msg *stakingtypes.MsgCreateValidator,
	cdc codec.Codec, broker rep.Broker, mapper tb.ToBroker) error {
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

	_ = desc

	// TODO: save to mongo?
	// TODO: test it
	if err = PublishValidatorsData(ctx, []types.StakingValidator{validator}, broker, mapper); err != nil {
		return err
	}

	// TODO: save to mongo?
	// TODO: test it
	// Save the first self-delegation
	if err = broker.PublishDelegation(ctx, mapper.MapDelegation(types.NewDelegation(
		msg.DelegatorAddress,
		msg.ValidatorAddress,
		msg.Value,
		height,
	))); err != nil {
		return err
	}

	// TODO:!!!!!
	// Save the description
	// err = db.SaveValidatorDescription(desc)
	// if err != nil {
	//	return err
	// }
	//

	// Save the commission
	// err = db.SaveValidatorCommission(types.NewValidatorCommission(
	//	msg.ValidatorAddress,
	//	&msg.Commission.Rate,
	//	&msg.MinSelfDelegation,
	//	height,
	// ))
	return err
}

// StoreDelegationFromMessage handles a MsgDelegate and saves the delegation inside the database
func StoreDelegationFromMessage(ctx context.Context, tx *types.Tx, msg *stakingtypes.MsgDelegate,
	stakingClient stakingtypes.QueryClient, broker rep.Broker, mapper tb.ToBroker) error {

	header := grpcClient.GetHeightRequestHeader(tx.Height)
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

	// TODO: test it
	d := types.NewDelegation(
		res.DelegationResponse.Delegation.DelegatorAddress,
		res.DelegationResponse.Delegation.ValidatorAddress,
		res.DelegationResponse.Balance,
		tx.Height,
	)

	if err = broker.PublishDelegation(ctx, mapper.MapDelegation(d)); err != nil {
		return err
	}

	dm := types.NewDelegationMessage(
		res.DelegationResponse.Delegation.DelegatorAddress,
		res.DelegationResponse.Delegation.ValidatorAddress,
		tx.TxHash,
		res.DelegationResponse.Balance,
		tx.Height,
	)

	if err = broker.PublishDelegationMessage(ctx, mapper.MapDelegationMessage(dm)); err != nil {
		return err
	}

	return nil
}

// StoreRedelegationFromMessage handles a MsgBeginRedelegate by saving the redelegation inside the database,
// and returns the new redelegation instance
func StoreRedelegationFromMessage(ctx context.Context, tx *types.Tx, index int, msg *stakingtypes.MsgBeginRedelegate,
	broker rep.Broker, mapper tb.ToBroker) (*types.Redelegation, error) {
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

	redelegationMessage := types.NewRedelegationMessage(
		msg.DelegatorAddress,
		msg.ValidatorSrcAddress,
		msg.ValidatorDstAddress,
		tx.TxHash,
		types.NewCoinFromCdk(msg.Amount),
		completionTime,
		tx.Height)

	// TODO: save to mongo?
	// TODO: test it
	err = broker.PublishRedelegation(ctx, mapper.MapRedelegation(redelegation))
	if err != nil {
		return nil, err
	}

	// TODO: test it
	err = broker.PublishRedelegationMessage(ctx, mapper.MapRedelegationMessage(redelegationMessage))
	if err != nil {
		return nil, err
	}

	return &redelegation, err
}

// StoreUnbondingDelegationFromMessage handles a MsgUndelegate storing the new unbonding delegation inside the database,
// and returns the new unbonding delegation instance
func StoreUnbondingDelegationFromMessage(ctx context.Context, tx *types.Tx, index int, msg *stakingtypes.MsgUndelegate,
	broker rep.Broker, mapper tb.ToBroker) (*types.UnbondingDelegation, error) {
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

	unbDelegation := types.NewUnbondingDelegation(
		msg.DelegatorAddress,
		msg.ValidatorAddress,
		msg.Amount,
		completionTime,
		tx.Height,
	)

	// TODO: test it
	err = broker.PublishUnbondingDelegation(ctx, mapper.MapUnbondingDelegation(unbDelegation))
	if err != nil {
		return nil, err
	}

	undDelegationMessage := types.NewUnbondingDelegationMessage(
		msg.DelegatorAddress,
		msg.ValidatorAddress,
		tx.TxHash,
		types.NewCoinFromCdk(msg.Amount),
		completionTime,
		tx.Height,
	)

	// TODO: test it
	err = broker.PublishUnbondingDelegationMessage(ctx, mapper.MapUnbondingDelegationMessage(undDelegationMessage))
	if err != nil {
		return nil, err
	}

	return &unbDelegation, err
}
