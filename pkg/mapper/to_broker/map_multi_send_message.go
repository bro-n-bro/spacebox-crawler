package to_broker

import (
	"github.com/hexy-dev/spacebox/broker/model"

	"github.com/hexy-dev/spacebox-crawler/types"
)

func (tb ToBroker) MapMsgMultiSend(msgSend types.MsgMultiSend) model.MultiSendMessage {
	return model.MultiSendMessage{
		Height:      msgSend.Height,
		AddressFrom: msgSend.AddressFrom,
		AddressTo:   msgSend.AddressTo,
		TxHash:      msgSend.TxHash,
		Coins:       tb.MapCoins(msgSend.Coins),
	}
}
