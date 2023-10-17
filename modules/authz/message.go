package authz

import (
	"context"

	codec "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"

	"github.com/bro-n-bro/spacebox-crawler/modules/utils"
	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// HandleMessageRecursive implements types.RecursiveMessagesHandler.
// Handles authz types messages.
// For MsgExec message types returns slice of messages to be handled recursively.
func (m *Module) HandleMessageRecursive(
	ctx context.Context,
	index int,
	cosmosMsg sdk.Msg,
	tx *types.Tx,
) ([]*codec.Any, error) {

	switch msg := cosmosMsg.(type) {
	case *authztypes.MsgGrant:
		if err := m.broker.PublishGrantMessage(ctx, model.GrantMessage{
			Height:     tx.Height,
			MsgIndex:   int64(index),
			TxHash:     tx.TxHash,
			Granter:    msg.Granter,
			Grantee:    msg.Grantee,
			Expiration: utils.TimeFromPtr(msg.Grant.Expiration),
			MsgType:    typeUrlFromAnyPtr(msg.Grant.Authorization),
		}); err != nil {
			return nil, err
		}

		if err := m.findAndPublishAuthzGrants(ctx, msg.Granter, msg.Grantee, tx.Height); err != nil {
			return nil, err
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
			return nil, err
		}

		if err := m.findAndPublishAuthzGrants(ctx, msg.Granter, msg.Grantee, tx.Height); err != nil {
			return nil, err
		}
	case *authztypes.MsgExec:
		messages := make([][]byte, 0, len(msg.Msgs))
		for _, message := range msg.Msgs {
			bytes, err := m.cdc.MarshalJSON(message)
			if err != nil {
				return nil, err
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
			return nil, err
		}

		return msg.Msgs, nil
	}

	return nil, nil
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
		if len(respPb.Grants) > 0 {
			for _, grant := range respPb.Grants {
				if err = m.broker.PublishAuthzGrant(ctx, model.AuthzGrant{
					Height:         height,
					GranterAddress: granter,
					GranteeAddress: grantee,
					Expiration:     utils.TimeFromPtr(grant.Expiration),
					MsgType:        typeUrlFromAnyPtr(grant.Authorization),
				}); err != nil {
					m.log.Err(err).Int64("height", height).Msg("error while publishing authz grant")
					return err
				}
			}
		} else {
			if err = m.broker.PublishAuthzGrant(ctx, model.AuthzGrant{
				Height:         height,
				GranterAddress: granter,
				GranteeAddress: grantee,
			}); err != nil {
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

// typeUrlFromAnyPtr returns typeUrl from *codec.Any. If any is nil, returns "".
// nolint: stylecheck
func typeUrlFromAnyPtr(any *codec.Any) string {
	if any == nil {
		return ""
	}

	return any.TypeUrl
}
