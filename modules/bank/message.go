package bank

import (
	"context"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"bro-n-bro-osmosis/types"
)

func (m *Module) HandleMessage(ctx context.Context, _ int, cosmosMsg sdk.Msg, tx *types.Tx) error {
	addresses, err := m.parser(m.cdc, cosmosMsg)
	if err != nil {
		m.log.Error().Err(err).Msg("HandleMessage getAddresses error:")
		return nil
	}
	//err = m.broker.PublishBank(ctx)

	switch msg := cosmosMsg.(type) {
	// todo: 	case *banktypes.MsgMultiSend:
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
