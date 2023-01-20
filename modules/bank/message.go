package bank

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/bro-n-bro/spacebox-crawler/modules/utils"
	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

func (m *Module) HandleMessage(ctx context.Context, index int, cosmosMsg sdk.Msg, tx *types.Tx) error {
	switch msg := cosmosMsg.(type) {
	case *banktypes.MsgMultiSend:
	// 	// todo: collect input/output and coins
	// 	// todo: think about how to collect total amount from outputs
	// 	if len(msg.Inputs) > 0 {
	// 		addressFrom := msg.Inputs[0].Address
	//
	// 		for _, to := range msg.Outputs {
	// 			if err := m.broker.PublishMultiSendMessage(ctx, model.MultiSendMessage{
	// 				Coins:       m.tbM.MapCoins(types.NewCoinsFromCdk(to.Coins)),
	// 				AddressFrom: addressFrom,
	// 				// AddressTo:   to.Address,
	// 				TxHash:   tx.TxHash,
	// 				Height:   tx.Height,
	// 				MsgIndex: int64(index),
	// 			}); err != nil {
	// 				return err
	// 			}
	// 		}
	// 	}
	case *banktypes.MsgSend:
		// TODO: test it
		if err := m.broker.PublishSendMessage(ctx, model.SendMessage{
			Coins:       m.tbM.MapCoins(types.NewCoinsFromCdk(msg.Amount)),
			AddressFrom: msg.FromAddress,
			AddressTo:   msg.ToAddress,
			TxHash:      tx.TxHash,
			Height:      tx.Height,
			MsgIndex:    int64(index),
		}); err != nil {
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
