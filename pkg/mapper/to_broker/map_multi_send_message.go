package to_broker

import (
	"github.com/hexy-dev/spacebox/broker/model"

	"bro-n-bro-osmosis/types"
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