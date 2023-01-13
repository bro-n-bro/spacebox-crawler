package utils

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authttypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/hexy-dev/spacebox-crawler/types"
)

// GetGenesisAccounts parses the given appState and returns the genesis accounts
func GetGenesisAccounts(appState map[string]json.RawMessage, cdc codec.Codec) ([]types.Account, error) {
	var authState authttypes.GenesisState
	if err := cdc.UnmarshalJSON(appState[authttypes.ModuleName], &authState); err != nil {
		return nil, err
	}

	// Store the accounts
	accounts := make([]types.Account, 0)

	for _, account := range authState.Accounts {
		var accountI authttypes.AccountI
		if err := cdc.UnpackAny(account, &accountI); err != nil {
			return nil, err
		}

		accounts = append(accounts, types.NewAccount(accountI.GetAddress().String(), 0))
	}

	return accounts, nil
}

// FilterNonAccountAddresses returns a slice containing only account addresses.
func FilterNonAccountAddresses(addresses []string) []string {
	// Filter using only the account addresses as the MessageAddressesParser might return also validator addresses
	accountAddresses := make([]string, 0)

	for _, address := range addresses {
		if _, err := sdk.AccAddressFromBech32(address); err == nil { // needs correct addresses only
			accountAddresses = append(accountAddresses, address)
		}
	}

	return accountAddresses
}
