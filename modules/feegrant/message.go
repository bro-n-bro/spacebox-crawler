package feegrant

import (
	"context"
	"time"

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

		var (
			allowance  feegranttypes.FeeAllowanceI
			expiration time.Time
		)
		if err := m.cdc.UnpackAny(msg.Allowance, &allowance); err != nil {
			return err
		}

		ex, err := allowance.ExpiresAt()
		if err != nil {
			return err
		}

		if ex != nil {
			expiration = *ex
		}

		data, err := m.cdc.MarshalJSON(msg.Allowance)
		if err != nil {
			return err
		}

		if err = m.broker.PublishGrantAllowanceMessage(ctx, model.GrantAllowanceMessage{
			Height:     tx.Height,
			MsgIndex:   int64(index),
			TxHash:     tx.TxHash,
			Granter:    msg.Granter,
			Grantee:    msg.Grantee,
			Expiration: expiration,
			Allowance:  data,
		}); err != nil {
			m.log.Err(err).Int64("height", tx.Height).Msg("error while publishing grant allowance message")
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
	}

	return nil
}
