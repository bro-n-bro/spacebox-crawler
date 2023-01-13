package staking

import (
	"context"
	"time"

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
		return stakingutils.StoreValidatorFromMsgCreateValidator(ctx, tx.Height, msg, m.cdc, m.tbM, m.broker)

	// TODO: does it needed?
	// case *stakingtypes.MsgEditValidator:
	//	return handleEditValidator(tx.Height, msg)

	case *stakingtypes.MsgDelegate:
		return stakingutils.StoreDelegationFromMessage(ctx, tx, msg, m.client.StakingQueryClient, m.tbM, m.broker)

	case *stakingtypes.MsgBeginRedelegate:
		return handleMsgBeginRedelegate(ctx, tx, index, m.tbM, m.broker, msg, m.client.StakingQueryClient)

	case *stakingtypes.MsgUndelegate:
		return handleMsgUndelegate(ctx, tx, index, m.tbM, m.broker, msg, m.client.StakingQueryClient)
	}

	return nil
}

// handleMsgBeginRedelegate handles and publishes a MsgBeginRedelegate data to broker
func handleMsgBeginRedelegate(ctx context.Context, tx *types.Tx, index int, mapper tb.ToBroker, broker broker,
	msg *stakingtypes.MsgBeginRedelegate, client stakingtypes.QueryClient) error {

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

	// TODO: save to mongo?
	// TODO: test it
	if err = broker.PublishRedelegation(ctx, model.Redelegation{
		Height:              tx.Height,
		DelegatorAddress:    msg.DelegatorAddress,
		SrcValidatorAddress: msg.ValidatorSrcAddress,
		DstValidatorAddress: msg.ValidatorDstAddress,
		Coin:                mapper.MapCoin(types.NewCoinFromCdk(msg.Amount)),
		CompletionTime:      completionTime,
	}); err != nil {
		return err
	}

	// TODO: test it
	if err = broker.PublishRedelegationMessage(ctx, model.RedelegationMessage{
		Redelegation: model.Redelegation{
			Height:              tx.Height,
			DelegatorAddress:    msg.DelegatorAddress,
			SrcValidatorAddress: msg.ValidatorSrcAddress,
			DstValidatorAddress: msg.ValidatorDstAddress,
			Coin:                mapper.MapCoin(types.NewCoinFromCdk(msg.Amount)),
			CompletionTime:      completionTime,
		},
		TxHash: tx.TxHash,
	}); err != nil {
		return err
	}

	// Update the current delegations
	return stakingutils.UpdateDelegationsAndReplaceExisting(ctx, tx.Height, msg.DelegatorAddress, client, mapper, broker)
}

// handleMsgUndelegate handles and publishes a MsgUndelegate data to broker
func handleMsgUndelegate(ctx context.Context, tx *types.Tx, index int, mapper tb.ToBroker, broker broker,
	msg *stakingtypes.MsgUndelegate, stakingClient stakingtypes.QueryClient) error {

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

	// TODO: test it
	if err = broker.PublishUnbondingDelegation(ctx, model.UnbondingDelegation{
		Height:              tx.Height,
		DelegatorAddress:    msg.DelegatorAddress,
		ValidatorAddress:    msg.ValidatorAddress,
		Coin:                mapper.MapCoin(types.NewCoinFromCdk(msg.Amount)),
		CompletionTimestamp: completionTime,
	}); err != nil {
		return err
	}

	// TODO: test it
	if err = broker.PublishUnbondingDelegationMessage(ctx, model.UnbondingDelegationMessage{
		UnbondingDelegation: model.UnbondingDelegation{
			Height:              tx.Height,
			DelegatorAddress:    msg.DelegatorAddress,
			ValidatorAddress:    msg.ValidatorAddress,
			Coin:                mapper.MapCoin(types.NewCoinFromCdk(msg.Amount)),
			CompletionTimestamp: completionTime,
		},
		TxHash: tx.TxHash,
	}); err != nil {
		return err
	}

	// Update the current delegations
	return stakingutils.UpdateDelegationsAndReplaceExisting(ctx, tx.Height, msg.DelegatorAddress, stakingClient,
		mapper, broker)
}
