package authz

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
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

		if err := m.findAndPublishAuthzGrants(ctx, msg.Granter, msg.Grantee, tx.Height); err != nil {
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

		if err := m.findAndPublishAuthzGrants(ctx, msg.Granter, msg.Grantee, tx.Height); err != nil {
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

// findAndPublishAuthzGrants finds all authz grants for the given granter and grantee and publishes data to the broker.
func (m *Module) findAndPublishAuthzGrants(ctx context.Context, granter, grantee string, height int64) error {
	var nextKey []byte

	for {
		respPb, err := m.client.AuthzQueryClient.Grants(
			ctx,
			&authztypes.QueryGrantsRequest{
				Granter: granter,
				Grantee: grantee,
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 150,
				},
			},
		)
		if err != nil {
			return err
		}

		nextKey = respPb.Pagination.NextKey
		for _, grant := range respPb.Grants {
			ag := model.AuthzGrant{
				Height:         height,
				GranterAddress: granter,
				GranteeAddress: grantee,
			}
			if grant.Expiration != nil {
				ag.Expiration = *grant.Expiration
			}

			if grant.Authorization != nil {
				ag.MsgType = grant.Authorization.TypeUrl
			}

			if err = m.broker.PublishAuthzGrant(ctx, ag); err != nil {
				m.log.Err(err).Int64("height", height).Msg("error while publishing authz grant")
				return err
			}
		}

		if len(respPb.Pagination.NextKey) == 0 {
			break
		}
	}

	return nil
}
