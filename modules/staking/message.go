package staking

import (
	"context"
	"time"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

func (m *Module) HandleMessage(ctx context.Context, index int, cosmosMsg sdk.Msg, tx *types.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch msg := cosmosMsg.(type) {
	case *stakingtypes.MsgCreateValidator:
		return m.handleMsgCreateValidator(ctx, tx.Height, tx.TxHash, index, msg)
	case *stakingtypes.MsgEditValidator:
		return m.handleEditValidator(ctx, tx.Height, tx.TxHash, index, msg)
	case *stakingtypes.MsgDelegate:
		return m.handleMsgDelegate(ctx, tx, msg, index)
	case *stakingtypes.MsgBeginRedelegate:
		return m.handleMsgBeginRedelegate(ctx, tx, index, msg)
	case *stakingtypes.MsgUndelegate:
		return m.handleMsgUndelegate(ctx, tx, index, msg)
	}

	return nil
}

// handleMsgCreateValidator handles MsgCreateValidator and publishes model.Validator, model.ValidatorDescription,
// model.Account, model.ValidatorInfo and model.Delegation messages to broker.
func (m *Module) handleMsgCreateValidator(
	ctx context.Context,
	height int64,
	hash string,
	index int,
	msg *stakingtypes.MsgCreateValidator,
) error {

	var pubKey cryptotypes.PubKey
	if err := m.cdc.UnpackAny(msg.Pubkey, &pubKey); err != nil {
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

	validator, err := convertValidator(m.cdc, stakingValidator, height)
	if err != nil {
		return err
	}
	commissionRate, err := stakingValidator.Commission.Rate.Float64()
	if err != nil {
		return err
	}
	if err = m.broker.PublishCreateValidatorMessage(ctx, model.CreateValidatorMessage{
		Height:           height,
		TxHash:           hash,
		MsgIndex:         int64(index),
		DelegatorAddress: msg.DelegatorAddress,
		ValidatorAddress: msg.ValidatorAddress,
		Description: model.ValidatorMessageDescription{
			Moniker:         msg.Description.Moniker,
			Identity:        msg.Description.Identity,
			Website:         msg.Description.Website,
			SecurityContact: msg.Description.SecurityContact,
			Details:         msg.Description.Details,
		},
		CommissionRates:   commissionRate,
		MinSelfDelegation: stakingValidator.GetMinSelfDelegation().Int64(),
		Pubkey:            pubKey.String(),
	}); err != nil {
		return err
	}

	if err = m.broker.PublishValidatorDescription(ctx, model.ValidatorDescription{
		OperatorAddress: msg.ValidatorAddress,
		Moniker:         msg.Description.Moniker,
		Identity:        msg.Description.Identity,
		Website:         msg.Description.Website,
		SecurityContact: msg.Description.SecurityContact,
		Details:         msg.Description.Details,
		AvatarURL:       "", // TODO
		Height:          height,
	}); err != nil {
		return err
	}

	// TODO: save to mongo?
	if err = m.PublishValidatorsData(ctx, []types.StakingValidator{validator}); err != nil {
		return err
	}

	// TODO: save to mongo?
	// Save the first self-delegation
	if err = m.broker.PublishDelegation(ctx, model.Delegation{
		OperatorAddress:  msg.ValidatorAddress,
		DelegatorAddress: msg.DelegatorAddress,
		Height:           height,
		Coin:             m.tbM.MapCoin(types.NewCoinFromCdk(msg.Value)),
	}); err != nil {
		return err
	}

	// FIXME: does it needed?
	// Save the description
	// err = broker.PublishValidatorDescription(desc)
	// if err != nil {
	//	return err
	// }
	//

	// Save the commission
	// err = broker.publishValidatorCommission(types.NewValidatorCommission(
	//	msg.ValidatorAddress,
	//	&msg.Commission.Rate,
	//	&msg.MinSelfDelegation,
	//	height,
	// ))
	return err
}

// handleMsgBeginRedelegate handles and publishes a MsgBeginRedelegate data to broker.
func (m *Module) handleMsgBeginRedelegate(ctx context.Context, tx *types.Tx, index int,
	msg *stakingtypes.MsgBeginRedelegate) error {

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
	if err = m.broker.PublishRedelegation(ctx, model.Redelegation{
		Height:              tx.Height,
		DelegatorAddress:    msg.DelegatorAddress,
		SrcValidatorAddress: msg.ValidatorSrcAddress,
		DstValidatorAddress: msg.ValidatorDstAddress,
		Coin:                m.tbM.MapCoin(types.NewCoinFromCdk(msg.Amount)),
		CompletionTime:      completionTime,
	}); err != nil {
		return err
	}

	// TODO: test it
	if err = m.broker.PublishRedelegationMessage(ctx, model.RedelegationMessage{
		Redelegation: model.Redelegation{
			Height:              tx.Height,
			DelegatorAddress:    msg.DelegatorAddress,
			SrcValidatorAddress: msg.ValidatorSrcAddress,
			DstValidatorAddress: msg.ValidatorDstAddress,
			Coin:                m.tbM.MapCoin(types.NewCoinFromCdk(msg.Amount)),
			CompletionTime:      completionTime,
		},
		TxHash:   tx.TxHash,
		MsgIndex: int64(index),
	}); err != nil {
		return err
	}

	// Update the current delegations
	return m.updateDelegationsAndReplaceExisting(ctx, tx.Height, msg.DelegatorAddress)
}

// handleMsgUndelegate handles MsgUndelegate and publishes data to broker.
func (m *Module) handleMsgUndelegate(ctx context.Context, tx *types.Tx, index int,
	msg *stakingtypes.MsgUndelegate) error {

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
	if err = m.broker.PublishUnbondingDelegation(ctx, model.UnbondingDelegation{
		Height:              tx.Height,
		DelegatorAddress:    msg.DelegatorAddress,
		ValidatorAddress:    msg.ValidatorAddress,
		Coin:                m.tbM.MapCoin(types.NewCoinFromCdk(msg.Amount)),
		CompletionTimestamp: completionTime,
	}); err != nil {
		return err
	}

	// TODO: test it
	if err = m.broker.PublishUnbondingDelegationMessage(ctx, model.UnbondingDelegationMessage{
		UnbondingDelegation: model.UnbondingDelegation{
			Height:              tx.Height,
			DelegatorAddress:    msg.DelegatorAddress,
			ValidatorAddress:    msg.ValidatorAddress,
			Coin:                m.tbM.MapCoin(types.NewCoinFromCdk(msg.Amount)),
			CompletionTimestamp: completionTime,
		},
		TxHash:   tx.TxHash,
		MsgIndex: int64(index),
	}); err != nil {
		return err
	}

	// Update the current delegations
	return m.updateDelegationsAndReplaceExisting(ctx, tx.Height, msg.DelegatorAddress)
}

// handleMsgDelegate handles a MsgDelegate and publish the delegation to broker.
func (m *Module) handleMsgDelegate(ctx context.Context, tx *types.Tx, msg *stakingtypes.MsgDelegate, index int) error {
	// TODO: test it
	if err := m.broker.PublishDelegationMessage(ctx, model.DelegationMessage{
		Delegation: model.Delegation{
			OperatorAddress:  msg.ValidatorAddress,
			DelegatorAddress: msg.DelegatorAddress,
			Coin:             m.tbM.MapCoin(types.NewCoinFromCdk(msg.Amount)),
			Height:           tx.Height,
		},
		TxHash:   tx.TxHash,
		MsgIndex: int64(index),
	}); err != nil {
		return err
	}

	header := grpcClient.GetHeightRequestHeader(tx.Height)

	respPb, err := m.client.StakingQueryClient.Delegation(
		ctx,
		&stakingtypes.QueryDelegationRequest{
			DelegatorAddr: msg.DelegatorAddress,
			ValidatorAddr: msg.ValidatorAddress,
		},
		header,
	)
	if err != nil {
		s, ok := status.FromError(err)
		if !ok || s.Code() != codes.NotFound {
			return err
		}
	}

	var coin model.Coin
	if err == nil {
		coin = m.tbM.MapCoin(types.NewCoinFromCdk(respPb.DelegationResponse.Balance))
	}

	// TODO: test it
	if err = m.broker.PublishDelegation(ctx, model.Delegation{
		OperatorAddress:  msg.ValidatorAddress,
		DelegatorAddress: msg.DelegatorAddress,
		Height:           tx.Height,
		Coin:             coin,
	}); err != nil {
		return err
	}

	return nil
}

// handleEditValidator handles MsgEditValidator and publishes model.ValidatorDescription to broker.
func (m *Module) handleEditValidator(ctx context.Context, height int64, hash string, index int, msg *stakingtypes.MsgEditValidator) error {
	if err := m.broker.PublishEditValidatorMessage(ctx, model.EditValidatorMessage{
		Height: height,
		Hash:   hash,
		Index:  int64(index),
		Description: model.ValidatorMessageDescription{
			Moniker:         msg.Description.Moniker,
			Identity:        msg.Description.Identity,
			Website:         msg.Description.Website,
			SecurityContact: msg.Description.SecurityContact,
			Details:         msg.Description.Details,
		},
	}); err != nil {
		return err
	}
	if err := m.broker.PublishValidatorDescription(ctx, model.ValidatorDescription{
		OperatorAddress: msg.ValidatorAddress,
		Moniker:         msg.Description.Moniker,
		Identity:        msg.Description.Identity,
		Website:         msg.Description.Website,
		SecurityContact: msg.Description.SecurityContact,
		Details:         msg.Description.Details,
		AvatarURL:       "", // TODO:
		Height:          height,
	}); err != nil {
		return err
	}
	return nil
}

// UpdateDelegationsAndReplaceExisting updates the delegations of the given delegator by querying them at the
// required height, and then publishes them to the broker by replacing all existing ones.
func (m *Module) updateDelegationsAndReplaceExisting(
	ctx context.Context,
	height int64,
	delegator string) error {
	// TODO:
	// Remove existing delegations
	// if err := broker.DeleteDelegatorDelegations(delegator); err != nil {
	//	return err
	// }

	// Get the delegations
	respPb, err := m.client.StakingQueryClient.DelegatorDelegations(
		ctx,
		&stakingtypes.QueryDelegatorDelegationsRequest{
			DelegatorAddr: delegator,
		},
	)
	if err != nil {
		return err
	}

	for _, delegation := range respPb.DelegationResponses {
		// TODO: test IT
		if err = m.broker.PublishDelegation(ctx, model.Delegation{
			OperatorAddress:  delegation.Delegation.ValidatorAddress,
			DelegatorAddress: delegation.Delegation.DelegatorAddress,
			Height:           height,
			Coin:             m.tbM.MapCoin(types.NewCoinFromCdk(delegation.Balance)),
		}); err != nil {
			return err
		}
	}

	return err
}
