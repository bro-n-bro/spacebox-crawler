package tobroker

import (
	"github.com/hexy-dev/spacebox/broker/model"
)

// TODO: use constructor
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
