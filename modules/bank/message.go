package bank

import (
	"context"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/hexy-dev/spacebox-crawler/types"
)

func (m *Module) HandleMessage(ctx context.Context, index int, cosmosMsg sdk.Msg, tx *types.Tx) error {
	addresses, err := m.parser(m.cdc, cosmosMsg)
	if err != nil {
		m.log.Error().Err(err).Msg("HandleMessage getAddresses error:")
		return nil
	}

	switch msg := cosmosMsg.(type) {
	case *banktypes.MsgMultiSend:
		// todo: collect input/output and coins
		// todo: think about how to collect total amount from outputs
		msgMultiSend := types.NewMsgMultiSend(nil, tx.Height, "", "", tx.TxHash)

		if err = m.broker.PublishMultiSendMessage(ctx, m.tbM.MapMsgMultiSend(msgMultiSend)); err != nil {
			return err
		}
	case *banktypes.MsgSend:
		// TODO: test it
		msgSend := types.NewMsgSend(types.NewCoinsFromCdk(msg.Amount), tx.Height, msg.FromAddress, msg.ToAddress, tx.TxHash)
		if err = m.broker.PublishSendMessage(ctx, m.tbM.MapMsgSend(msgSend)); err != nil {
			return err
		}
	}

	// todo: publish?
	_ = addresses

	return nil
}
