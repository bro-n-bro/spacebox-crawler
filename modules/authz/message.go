package authz

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"

	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

func (m *Module) HandleMessage(ctx context.Context, index int, cosmosMsg sdk.Msg, tx *types.Tx) error {

	switch msg := cosmosMsg.(type) {
	case *authztypes.MsgGrant:
		// m.client.AuthzQueryClient.Grants(ctx, &authztypes.MsgGrant{})
		var auth authztypes.Authorization
		if err := m.cdc.UnpackAny(msg.Grant.Authorization, &auth); err != nil {
			return err
		}
		var expiration time.Time
		if msg.Grant.Expiration != nil {
			expiration = *msg.Grant.Expiration
		}

		if err := m.broker.PublishGrantMessage(ctx, model.GrantMessage{
			Height:     tx.Height,
			MsgIndex:   int64(index),
			TxHash:     tx.TxHash,
			Granter:    msg.Granter,
			Grantee:    msg.Grantee,
			Expiration: expiration,
			MsgType:    auth.MsgTypeURL(),
		}); err != nil {
			return err
		}
	case *authztypes.MsgRevoke:
	case *authztypes.MsgExec:
	}
	return nil
}
