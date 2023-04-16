package feegrant

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	feegranttypes "github.com/cosmos/cosmos-sdk/x/feegrant"

	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

func (m *Module) HandleMessage(ctx context.Context, index int, cosmosMsg sdk.Msg, tx *types.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch msg := cosmosMsg.(type) {
	case *feegranttypes.MsgGrantAllowance:
		data, err := m.cdc.MarshalJSON(msg.Allowance)
		if err != nil {
			return err
		}

		if err = m.broker.PublishGrantAllowanceMessage(ctx, model.GrantAllowanceMessage{
			Height:    tx.Height,
			MsgIndex:  int64(index),
			TxHash:    tx.TxHash,
			Granter:   msg.Granter,
			Grantee:   msg.Grantee,
			Allowance: data,
		}); err != nil {
			m.log.Err(err).Int64("height", tx.Height).Msg("error while publishing grant allowance message")
			return err
		}

		if err = m.publishFeeAllowance(ctx, msg.Granter, msg.Grantee); err != nil {
			m.log.Err(err).
				Int64("height", tx.Height).
				Str("message", "MsgGrantAllowance").
				Msg("error while publishing fee allowance")
			return err
		}

	case *feegranttypes.MsgRevokeAllowance:
		if err := m.broker.PublishRevokeAllowanceMessage(ctx, model.RevokeAllowanceMessage{
			Height:   tx.Height,
			MsgIndex: int64(index),
			TxHash:   tx.TxHash,
			Granter:  msg.Granter,
			Grantee:  msg.Grantee,
		}); err != nil {
			m.log.Err(err).Int64("height", tx.Height).Msg("error while publishing grant allowance message")
			return err
		}

		if err := m.publishFeeAllowance(ctx, msg.Granter, msg.Grantee); err != nil {
			m.log.Err(err).
				Int64("height", tx.Height).
				Str("message", "MsgRevokeAllowance").
				Msg("error while publishing fee allowance")
			return err
		}
	}

	return nil
}

func (m *Module) publishFeeAllowance(ctx context.Context, granter, grantee string) error {
	respPb, err := m.client.FeegrantQueryClient.Allowance(ctx, &feegranttypes.QueryAllowanceRequest{
		Granter: granter,
		Grantee: grantee,
	})
	if err != nil {
		return err
	}

	allowanceBytes, err := m.cdc.MarshalJSON(respPb.Allowance)
	if err != nil {
		return err
	}

	return m.broker.PublishFeeAllowance(ctx, model.FeeAllowance{
		Granter:   granter,
		Grantee:   grantee,
		Allowance: allowanceBytes,
	})
}
