package to_broker

import (
	"bro-n-bro-osmosis/adapter/broker/model"
	"bro-n-bro-osmosis/types"
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
