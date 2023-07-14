package ibc

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"

	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

func (m *Module) HandleMessage(ctx context.Context, index int, cosmosMsg sdk.Msg, tx *types.Tx) error {
	switch msg := cosmosMsg.(type) {
	case *ibctransfertypes.MsgTransfer:
		return m.broker.PublishTransferMessage(ctx, model.TransferMessage{
			SourceChannel: msg.SourceChannel,
			Coin:          m.tbM.MapCoin(types.NewCoinFromCdk(msg.Token)),
			Sender:        msg.Sender,
			Receiver:      msg.Receiver,
			Height:        tx.Height,
			MsgIndex:      int64(index),
			TxHash:        tx.TxHash,
		})
	case *ibcchanneltypes.MsgAcknowledgement:
		return m.broker.PublishAcknowledgementMessage(ctx, model.AcknowledgementMessage{
			SourcePort:         msg.Packet.SourcePort,
			SourceChannel:      msg.Packet.SourceChannel,
			DestinationPort:    msg.Packet.DestinationPort,
			DestinationChannel: msg.Packet.DestinationChannel,
			Data:               msg.Packet.Data,
			TimeoutTimestamp:   msg.Packet.TimeoutTimestamp,
			Signer:             msg.Signer,
			Height:             tx.Height,
			MsgIndex:           int64(index),
			TxHash:             tx.TxHash,
		})
	case *ibcchanneltypes.MsgRecvPacket:
		return m.broker.PublishReceivePacketMessage(ctx, model.RecvPacketMessage{
			SourcePort:         msg.Packet.SourcePort,
			SourceChannel:      msg.Packet.SourceChannel,
			DestinationPort:    msg.Packet.DestinationPort,
			DestinationChannel: msg.Packet.DestinationChannel,
			Signer:             msg.Signer,
			Data:               msg.Packet.Data,
			TimeoutTimestamp:   msg.Packet.TimeoutTimestamp,
			Height:             tx.Height,
			MsgIndex:           int64(index),
			TxHash:             tx.TxHash,
		})
	}

	return nil
}
