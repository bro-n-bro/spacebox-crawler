package utils

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	authttypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"golang.org/x/exp/constraints"

	"github.com/hexy-dev/spacebox-crawler/types"
)

// GetGenesisAccounts parses the given appState and returns the genesis accounts
func GetGenesisAccounts(appState map[string]json.RawMessage, cdc codec.Codec) ([]types.Account, error) {
	var authState authttypes.GenesisState
	if err := cdc.UnmarshalJSON(appState[authttypes.ModuleName], &authState); err != nil {
		return nil, err
	}

	// Store the accounts
	accounts := make([]types.Account, len(authState.Accounts))
	for index, account := range authState.Accounts {
		var accountI authttypes.AccountI
		err := cdc.UnpackAny(account, &accountI)
		if err != nil {
			return nil, err
		}

		accounts[index] = types.NewAccount(accountI.GetAddress().String(), 0)
	}

	return accounts, nil
}

func GetAccounts(addresses []string, height int64) []types.Account {
	res := make([]types.Account, len(addresses))
	for i, addr := range addresses {
		res[i] = types.NewAccount(addr, height)
	}
	return res
}

func ContainAny[T constraints.Ordered](src []T, trg T) bool {
	for _, v := range src {
		if v == trg {
			return true
		}
	}
	return false
}
