package bank

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/hexy-dev/spacebox-crawler/modules/utils"
	"github.com/hexy-dev/spacebox-crawler/types"
	"github.com/hexy-dev/spacebox/broker/model"
)

func (m *Module) HandleMessage(ctx context.Context, _ int, cosmosMsg sdk.Msg, tx *types.Tx) error {
	switch msg := cosmosMsg.(type) {
	case *banktypes.MsgMultiSend:
		// todo: collect input/output and coins
		// todo: think about how to collect total amount from outputs
		if len(msg.Inputs) > 0 {
			addressFrom := msg.Inputs[0].Address
			for _, to := range msg.Outputs {
				if err := m.broker.PublishMultiSendMessage(ctx, model.NewMultiSendMessage(tx.Height, addressFrom,
					to.Address, tx.TxHash, m.tbM.MapCoins(types.NewCoinsFromCdk(to.Coins)))); err != nil {
					return err
				}
			}
		}
	case *banktypes.MsgSend:
		msgSend := model.NewSendMessage(tx.Height, msg.FromAddress, msg.ToAddress, tx.TxHash,
			m.tbM.MapCoins(types.NewCoinsFromCdk(msg.Amount)))

		// TODO: test it
		if err := m.broker.PublishSendMessage(ctx, msgSend); err != nil {
			return err
		}
	}

	addresses, err := m.parser(m.cdc, cosmosMsg)
	if err != nil {
		m.log.Error().Err(err).Msg("HandleMessage getAddresses error:")
		return nil
	}

	return m.updateBalance(ctx, utils.FilterNonAccountAddresses(addresses), tx.Height)
}
