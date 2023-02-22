package tobroker

import (
	"github.com/bro-n-bro/spacebox/broker/model"
)

func (tb ToBroker) MapMessage(txHash, msgType, signer string, index int, accounts []string, value []byte) model.Message {
	return model.Message{
		TransactionHash:           txHash,
		MsgIndex:                  int64(index),
		Type:                      msgType,
		InvolvedAccountsAddresses: accounts,
		Signer:                    signer,
		Value:                     value,
	}
}
