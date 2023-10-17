package ibc

import (
	"context"
	"strings"

	codec "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	jsoniter "github.com/json-iterator/go"

	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

type icaHostData struct {
	Data []byte `json:"data"`
}

// HandleMessageRecursive implements types.RecursiveMessagesHandler.
// Handles ibc types messages.
// For MsgTransfer message types returns slice of messages to be handled recursively.
func (m *Module) HandleMessageRecursive(
	ctx context.Context,
	index int,
	cosmosMsg sdk.Msg,
	tx *types.Tx,
) ([]*codec.Any, error) {

	switch msg := cosmosMsg.(type) {
	case *ibctransfertypes.MsgTransfer:
		return nil, m.broker.PublishTransferMessage(ctx, model.TransferMessage{
			SourceChannel: msg.SourceChannel,
			Coin:          m.tbM.MapCoin(types.NewCoinFromCdk(msg.Token)),
			Sender:        msg.Sender,
			Receiver:      msg.Receiver,
			Height:        tx.Height,
			MsgIndex:      int64(index),
			TxHash:        tx.TxHash,
		})
	case *ibcchanneltypes.MsgAcknowledgement:
		return nil, m.broker.PublishAcknowledgementMessage(ctx, model.AcknowledgementMessage{
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
		if err := m.broker.PublishReceivePacketMessage(ctx, model.RecvPacketMessage{
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
		}); err != nil {
			return nil, err
		}

		if msg.Packet.DestinationPort == "icahost" {
			hostData := icaHostData{}
			if err := jsoniter.Unmarshal(msg.Packet.Data, &hostData); err != nil {
				return nil, err
			}

			cosmosTx := icatypes.CosmosTx{}
			if err := m.cdc.Unmarshal(hostData.Data, &cosmosTx); err != nil {
				// skip unsupported messages
				if strings.HasPrefix(err.Error(), "no concrete type registered for type URL") {
					m.log.Warn().Err(err).Msg("error while unpacking message")
					return nil, nil
				}

				return nil, err
			}

			return cosmosTx.Messages, nil
		}
	}

	return nil, nil
}
