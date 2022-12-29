package to_broker

import (
	"github.com/hexy-dev/spacebox-crawler/types"
	"github.com/hexy-dev/spacebox/broker/model"
)

func (tb ToBroker) MapMsgSend(msgSend types.MsgSend) model.SendMessage {
	return model.SendMessage{
		Height:      msgSend.Height,
		AddressFrom: msgSend.AddressFrom,
		AddressTo:   msgSend.AddressTo,
		TxHash:      msgSend.TxHash,
		Coins:       tb.MapCoins(msgSend.Coins),
	}
}
