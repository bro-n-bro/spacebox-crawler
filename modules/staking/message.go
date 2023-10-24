package staking

import (
	"context"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/pkg/errors"
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
	case *stakingtypes.MsgCancelUnbondingDelegation:
		return m.handleMsgCancelUnbondingDelegation(ctx, tx, index, msg)
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
		Coin:             m.tbM.MapCoin(types.NewCoinFromSDK(msg.Value)),
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

	// try to find the completion time in event. It does not exist in IBC transactions
	completionTime := findCompletionTimeInEventOrZero(tx, index, stakingtypes.EventTypeRedelegate)

	var (
		redelegationsResp []stakingtypes.RedelegationResponse
		nextKey           []byte
	)

	for {
		respPb, err := m.client.StakingQueryClient.Redelegations(ctx, &stakingtypes.QueryRedelegationsRequest{
			DelegatorAddr: msg.DelegatorAddress,
			Pagination: &query.PageRequest{
				Key:        nextKey,
				Limit:      100,
				CountTotal: true,
			},
		})
		if err != nil {
			s, ok := status.FromError(err)
			if !ok {
				return err
			}

			if s.Code() != codes.NotFound {
				return err
			}

			goto Publish
		}

		// first iteration
		if len(nextKey) == 0 {
			redelegationsResp = make([]stakingtypes.RedelegationResponse, 0, respPb.Pagination.Total)
		}

		nextKey = respPb.Pagination.NextKey
		redelegationsResp = append(redelegationsResp, respPb.RedelegationResponses...)

		if len(respPb.Pagination.NextKey) == 0 {
			break
		}
	}

	for _, resp := range redelegationsResp {
		for _, entry := range resp.Entries {
			if entry.RedelegationEntry.CreationHeight == tx.Height {
				completionTime = entry.RedelegationEntry.CompletionTime
				continue // we will publish it in the publish section
			}

			if err := m.broker.PublishRedelegation(ctx, model.Redelegation{
				Height:              entry.RedelegationEntry.CreationHeight,
				DelegatorAddress:    resp.Redelegation.DelegatorAddress,
				SrcValidatorAddress: resp.Redelegation.ValidatorSrcAddress,
				DstValidatorAddress: resp.Redelegation.ValidatorDstAddress,
				Coin:                m.tbM.MapCoin(types.NewCoin(m.defaultDenom, float64(entry.Balance.BigInt().Int64()))), // nolint: lll
				CompletionTime:      entry.RedelegationEntry.CompletionTime,
			}); err != nil {
				return err
			}
		}
	}

Publish:
	// TODO: save to mongo?
	if err := m.broker.PublishRedelegation(ctx, model.Redelegation{
		Height:              tx.Height,
		DelegatorAddress:    msg.DelegatorAddress,
		SrcValidatorAddress: msg.ValidatorSrcAddress,
		DstValidatorAddress: msg.ValidatorDstAddress,
		Coin:                m.tbM.MapCoin(types.NewCoinFromSDK(msg.Amount)),
		CompletionTime:      completionTime,
	}); err != nil {
		return err
	}

	if err := m.broker.PublishRedelegationMessage(ctx, model.RedelegationMessage{
		Redelegation: model.Redelegation{
			Height:              tx.Height,
			DelegatorAddress:    msg.DelegatorAddress,
			SrcValidatorAddress: msg.ValidatorSrcAddress,
			DstValidatorAddress: msg.ValidatorDstAddress,
			Coin:                m.tbM.MapCoin(types.NewCoinFromSDK(msg.Amount)),
			CompletionTime:      completionTime,
		},
		TxHash:   tx.TxHash,
		MsgIndex: int64(index),
	}); err != nil {
		return err
	}

	// Update the current delegations
	if err := m.updateDelegations(ctx, tx.Height, msg.DelegatorAddress, msg.ValidatorSrcAddress); err != nil {
		return errors.Wrap(err, "update delegations")
	}

	// check the delegations for the current validator
	return m.updateOrDisableDelegation(ctx, msg.DelegatorAddress, msg.ValidatorDstAddress, tx.Height)
}

// handleMsgUndelegate handles MsgUndelegate and publishes data to broker.
func (m *Module) handleMsgUndelegate(ctx context.Context, tx *types.Tx, index int,
	msg *stakingtypes.MsgUndelegate) error {

	// try to find the completion time in event. It does not exist in IBC transactions
	completionTime := findCompletionTimeInEventOrZero(tx, index, stakingtypes.EventTypeUnbond)

	respPb, err := m.client.StakingQueryClient.UnbondingDelegation(ctx, &stakingtypes.QueryUnbondingDelegationRequest{
		DelegatorAddr: msg.DelegatorAddress,
		ValidatorAddr: msg.ValidatorAddress,
	})
	if err != nil {
		s, ok := status.FromError(err)
		if !ok {
			return err
		}

		if s.Code() != codes.NotFound {
			return err
		}

		goto PublishMessage
	}

	for _, entry := range respPb.Unbond.Entries {
		if entry.CreationHeight == tx.Height {
			completionTime = entry.CompletionTime
		}

		if err = m.broker.PublishUnbondingDelegation(ctx, model.UnbondingDelegation{
			Height:           entry.CreationHeight,
			DelegatorAddress: respPb.Unbond.DelegatorAddress,
			OperatorAddress:  respPb.Unbond.ValidatorAddress,
			Coin:             m.tbM.MapCoin(types.NewCoin(m.defaultDenom, float64(entry.Balance.BigInt().Int64()))), //nolint:lll
			CompletionTime:   entry.CompletionTime,
		}); err != nil {
			return err
		}
	}

PublishMessage:
	if err = m.broker.PublishUnbondingDelegationMessage(ctx, model.UnbondingDelegationMessage{
		UnbondingDelegation: model.UnbondingDelegation{
			Height:           tx.Height,
			DelegatorAddress: msg.DelegatorAddress,
			OperatorAddress:  msg.ValidatorAddress,
			Coin:             m.tbM.MapCoin(types.NewCoinFromSDK(msg.Amount)),
			CompletionTime:   completionTime,
		},
		TxHash:   tx.TxHash,
		MsgIndex: int64(index),
	}); err != nil {
		return err
	}

	// Update the current delegations
	return m.updateDelegations(ctx, tx.Height, msg.DelegatorAddress, msg.ValidatorAddress)
}

// handleMsgDelegate handles a MsgDelegate and publish the delegation to broker.
func (m *Module) handleMsgDelegate(ctx context.Context, tx *types.Tx, msg *stakingtypes.MsgDelegate, index int) error {
	if err := m.broker.PublishDelegationMessage(ctx, model.DelegationMessage{
		Delegation: model.Delegation{
			OperatorAddress:  msg.ValidatorAddress,
			DelegatorAddress: msg.DelegatorAddress,
			Coin:             m.tbM.MapCoin(types.NewCoinFromSDK(msg.Amount)),
			Height:           tx.Height,
		},
		TxHash:   tx.TxHash,
		MsgIndex: int64(index),
	}); err != nil {
		return err
	}

	if err := m.updateOrDisableDelegation(ctx, msg.DelegatorAddress, msg.ValidatorAddress, tx.Height); err != nil {
		return errors.Wrap(err, "update or disable delegation")
	}

	return nil
}

// handleEditValidator handles MsgEditValidator and publishes model.ValidatorDescription to broker.
func (m *Module) handleEditValidator(
	ctx context.Context,
	height int64,
	hash string,
	index int,
	msg *stakingtypes.MsgEditValidator,
) error {

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
		Height:          height,
	}); err != nil {
		return err
	}

	return nil
}

// updateDelegations updates the delegations of the given delegator by querying them at the
// required height, and then publishes them to the broker by replacing all existing ones.
//
// also checks the delegation with the current validator address,
// and publishes disabled delegation to the broker if it doesn't exist.
func (m *Module) updateDelegations(ctx context.Context, height int64, delegator, validator string) error {
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

	var delegationWithCurValidatorExists bool

	for _, delegation := range respPb.DelegationResponses {
		if delegation.Delegation.ValidatorAddress == validator {
			delegationWithCurValidatorExists = true
		}

		if err = m.broker.PublishDelegation(ctx, model.Delegation{
			OperatorAddress:  delegation.Delegation.ValidatorAddress,
			DelegatorAddress: delegation.Delegation.DelegatorAddress,
			Height:           height,
			Coin:             m.tbM.MapCoin(types.NewCoinFromSDK(delegation.Balance)),
		}); err != nil {
			return err
		}
	}

	if !delegationWithCurValidatorExists {
		if err = m.updateOrDisableDelegation(ctx, delegator, validator, height); err != nil {
			return errors.Wrap(err, "disable delegation")
		}
	}

	return nil
}

// updateOrDisableDelegation checks the delegation with the given validator address
// if it exists, publishes updated delegation to the broker, otherwise publishes disabled delegation.
func (m *Module) updateOrDisableDelegation(ctx context.Context, delegatorAddr, operatorAddr string, height int64) error {
	respPb, err := m.client.StakingQueryClient.Delegation(ctx,
		&stakingtypes.QueryDelegationRequest{
			DelegatorAddr: delegatorAddr,
			ValidatorAddr: operatorAddr,
		},
		grpcClient.GetHeightRequestHeader(height),
	)
	if err != nil {
		s, ok := status.FromError(err)
		if !ok || s.Code() != codes.NotFound {
			return err
		}

		return m.broker.PublishDisabledDelegation(ctx, model.Delegation{
			OperatorAddress:  operatorAddr,
			DelegatorAddress: delegatorAddr,
			Coin:             model.Coin{}, // zero coin
			Height:           height,
		})
	}

	if err = m.broker.PublishDelegation(ctx, model.Delegation{
		OperatorAddress:  operatorAddr,
		DelegatorAddress: delegatorAddr,
		Height:           height,
		Coin:             m.tbM.MapCoin(types.NewCoinFromSDK(respPb.DelegationResponse.Balance)),
	}); err != nil {
		return err
	}

	return nil
}

// handleMsgCancelUnbondingDelegation handles MsgCancelUnbondingDelegation
// and publishes model.CancelUnbondingDelegationMessage to broker.
func (m *Module) handleMsgCancelUnbondingDelegation(
	ctx context.Context,
	tx *types.Tx,
	index int,
	msg *stakingtypes.MsgCancelUnbondingDelegation,
) error {

	return m.broker.PublishCancelUnbondingDelegationMessage(ctx, model.CancelUnbondingDelegationMessage{
		Height:           tx.Height,
		TxHash:           tx.TxHash,
		MsgIndex:         int64(index),
		ValidatorAddress: msg.ValidatorAddress,
		DelegatorAddress: msg.DelegatorAddress,
	})
}
