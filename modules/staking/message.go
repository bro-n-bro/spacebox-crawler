package staking

import (
	"context"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	stakingutils "github.com/hexy-dev/spacebox-crawler/modules/staking/utils"
	tb "github.com/hexy-dev/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/hexy-dev/spacebox-crawler/types"
	"github.com/hexy-dev/spacebox/broker/model"
)

func (m *Module) HandleMessage(ctx context.Context, index int, cosmosMsg sdk.Msg, tx *types.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch msg := cosmosMsg.(type) {
	case *stakingtypes.MsgCreateValidator:
		return handleMsgCreateValidator(ctx, tx.Height, msg, m.cdc, m.tbM, m.broker)

	// TODO: does it needed?
	// case *stakingtypes.MsgEditValidator:
	//	return handleEditValidator(tx.Height, msg)

	case *stakingtypes.MsgDelegate:
		return stakingutils.StoreDelegationFromMessage(ctx, tx, msg, m.client.StakingQueryClient, m.tbM, m.broker)

	case *stakingtypes.MsgBeginRedelegate:
		return handleMsgBeginRedelegate(ctx, tx, index, msg, m.client.StakingQueryClient, m.tbM, m.broker)

	case *stakingtypes.MsgUndelegate:
		return handleMsgUndelegate(ctx, tx, index, msg, m.client.StakingQueryClient, m.tbM, m.broker)
	}

	return nil
}

// handleMsgCreateValidator handles properly a MsgCreateValidator instance by
// saving into the database all the data associated to such validator
func handleMsgCreateValidator(ctx context.Context, height int64, msg *stakingtypes.MsgCreateValidator, cdc codec.Codec,
	mapper tb.ToBroker, broker broker) error {

	return stakingutils.StoreValidatorFromMsgCreateValidator(ctx, height, msg, cdc, mapper, broker)
}

// handleMsgBeginRedelegate handles and publishes a MsgBeginRedelegate data to broker
func handleMsgBeginRedelegate(ctx context.Context, tx *types.Tx, index int, msg *stakingtypes.MsgBeginRedelegate,
	client stakingtypes.QueryClient, mapper tb.ToBroker, broker broker) error {

	event, err := tx.FindEventByType(index, stakingtypes.EventTypeRedelegate)
	if err != nil {
		return err
	}

	completionTimeStr, err := tx.FindAttributeByKey(event, stakingtypes.AttributeKeyCompletionTime)
	if err != nil {
		return err
	}

	completionTime, err := time.Parse(time.RFC3339, completionTimeStr)
	if err != nil {
		return err
	}

	redelegation := model.NewRedelegation(
		tx.Height,
		msg.DelegatorAddress,
		msg.ValidatorSrcAddress,
		msg.ValidatorDstAddress,
		mapper.MapCoin(types.NewCoinFromCdk(msg.Amount)),
		completionTime,
	)

	// TODO: save to mongo?
	// TODO: test it
	err = broker.PublishRedelegation(ctx, redelegation)
	if err != nil {
		return err
	}

	redelegationMessage := model.NewRedelegationMessage(
		tx.Height,
		msg.DelegatorAddress,
		msg.ValidatorSrcAddress,
		msg.ValidatorDstAddress,
		tx.TxHash,
		mapper.MapCoin(types.NewCoinFromCdk(msg.Amount)),
		completionTime,
	)

	// TODO: test it
	err = broker.PublishRedelegationMessage(ctx, redelegationMessage)
	if err != nil {
		return err
	}

	// Update the current delegations
	return stakingutils.UpdateDelegationsAndReplaceExisting(ctx, tx.Height, msg.DelegatorAddress, client, mapper, broker)
}

// handleMsgUndelegate handles and publishes a MsgUndelegate data to broker
func handleMsgUndelegate(ctx context.Context, tx *types.Tx, index int, msg *stakingtypes.MsgUndelegate,
	stakingClient stakingtypes.QueryClient, mapper tb.ToBroker, broker broker) error {

	event, err := tx.FindEventByType(index, stakingtypes.EventTypeUnbond)
	if err != nil {
		return err
	}

	completionTimeStr, err := tx.FindAttributeByKey(event, stakingtypes.AttributeKeyCompletionTime)
	if err != nil {
		return err
	}

	completionTime, err := time.Parse(time.RFC3339, completionTimeStr)
	if err != nil {
		return err
	}

	unbDelegation := model.NewUnbondingDelegation(
		tx.Height,
		msg.DelegatorAddress,
		msg.ValidatorAddress,
		mapper.MapCoin(types.NewCoinFromCdk(msg.Amount)),
		completionTime,
	)

	// TODO: test it
	err = broker.PublishUnbondingDelegation(ctx, unbDelegation)
	if err != nil {
		return err
	}

	undDelegationMessage := model.NewUnbondingDelegationMessage(
		tx.Height,
		msg.DelegatorAddress,
		msg.ValidatorAddress,
		tx.TxHash,
		mapper.MapCoin(types.NewCoinFromCdk(msg.Amount)),
		completionTime,
	)

	// TODO: test it
	err = broker.PublishUnbondingDelegationMessage(ctx, undDelegationMessage)
	if err != nil {
		return err
	}

	// Update the current delegations
	return stakingutils.UpdateDelegationsAndReplaceExisting(ctx, tx.Height, msg.DelegatorAddress, stakingClient, mapper, broker)
}
