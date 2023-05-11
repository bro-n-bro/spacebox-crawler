package slashing

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	slashtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"

	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

func (m *Module) HandleMessage(ctx context.Context, index int, cosmosMsg sdk.Msg, tx *types.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch msg := cosmosMsg.(type) {
	case *slashtypes.MsgUnjail:
		return m.handleMsgUnjail(ctx, tx, index, msg)
	default:
		return nil
	}
}

func (m *Module) handleMsgUnjail(ctx context.Context, tx *types.Tx, index int, msg *slashtypes.MsgUnjail) error {
	return m.broker.PublishUnjailMessage(ctx, model.UnjailMessage{
		Height:        tx.Height,
		Hash:          tx.TxHash,
		Index:         int64(index),
		ValidatorAddr: msg.ValidatorAddr,
	})
}
