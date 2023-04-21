package authz

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"

	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// HandleMessage implements types.MessageHandler.
// Handles authz types messages.
func (m *Module) HandleMessage(ctx context.Context, index int, cosmosMsg sdk.Msg, tx *types.Tx) error {
	switch msg := cosmosMsg.(type) {
	case *authztypes.MsgGrant:
		var (
			expiration time.Time
			msgType    string
		)
		if msg.Grant.Expiration != nil {
			expiration = *msg.Grant.Expiration
		}
		if msg.Grant.Authorization != nil {
			msgType = msg.Grant.Authorization.TypeUrl
		}

		if err := m.broker.PublishGrantMessage(ctx, model.GrantMessage{
			Height:     tx.Height,
			MsgIndex:   int64(index),
			TxHash:     tx.TxHash,
			Granter:    msg.Granter,
			Grantee:    msg.Grantee,
			Expiration: expiration,
			MsgType:    msgType,
		}); err != nil {
			return err
		}

	case *authztypes.MsgRevoke:
		if err := m.broker.PublishRevokeMessage(ctx, model.RevokeMessage{
			Height:   tx.Height,
			MsgIndex: int64(index),
			TxHash:   tx.TxHash,
			Granter:  msg.Granter,
			Grantee:  msg.Grantee,
			MsgType:  msg.MsgTypeUrl,
		}); err != nil {
			return err
		}

	case *authztypes.MsgExec:
		messages := make([][]byte, 0, len(msg.Msgs))
		for _, message := range msg.Msgs {
			bytes, err := m.cdc.MarshalJSON(message)
			if err != nil {
				return err
			}
			messages = append(messages, bytes)
		}
		if err := m.broker.PublishExecMessage(ctx, model.ExecMessage{
			Height:   tx.Height,
			MsgIndex: int64(index),
			TxHash:   tx.TxHash,
			Grantee:  msg.Grantee,
			Msgs:     messages,
		}); err != nil {
			return err
		}
	}

	return nil
}
