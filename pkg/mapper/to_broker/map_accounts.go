package to_broker

import (
	"github.com/hexy-dev/spacebox/broker/model"

	"github.com/hexy-dev/spacebox-crawler/types"
)

func (tb ToBroker) MapAccount(account types.Account) model.Account {
	return model.Account{
		Address: account.Address,
		Height:  account.Height,
	}
}

func (tb ToBroker) MapAccounts(accounts []types.Account) []model.Account {
	res := make([]model.Account, len(accounts))
	for i, acc := range accounts {
		res[i] = tb.MapAccount(acc)
	}
	return res
}
