package to_broker

import (
	"bro-n-bro-osmosis/adapter/broker/model"
)

func (tb ToBroker) MapMessage(txHash, msgType, signer string, index int, accounts []string, value []byte) model.Message {
	return model.Message{
		TransactionHash:           txHash,
		Index:                     index,
		Type:                      msgType,
		InvolvedAccountsAddresses: accounts,
		Signer:                    signer,
		Value:                     value,
	}
}
