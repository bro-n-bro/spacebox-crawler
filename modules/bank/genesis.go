package bank

import (
	"context"
	"encoding/json"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"bro-n-bro-osmosis/modules/utils"
	"bro-n-bro-osmosis/types"
)

// HandleGenesis handles the genesis state of the x/bank module in order to store the initial values
// of the different account balances.
func (m *Module) HandleGenesis(_ context.Context, doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	var bankState banktypes.GenesisState
	if err := m.cdc.UnmarshalJSON(appState[banktypes.ModuleName], &bankState); err != nil {
		return err
	}

	// Store the balances
	accounts, err := utils.GetGenesisAccounts(appState, m.cdc)
	if err != nil {
		return err
	}
	accountsMap := getAccountsMap(accounts)

	var balances []types.AccountBalance
	for _, balance := range bankState.Balances {
		_, ok := accountsMap[balance.Address]
		if !ok {
			continue
		}

		balances = append(balances, types.NewAccountBalance(balance.Address, balance.Coins, doc.InitialHeight))
	}

	// TODO:
	_ = balances
	return nil
}

func getAccountsMap(accounts []types.Account) map[string]struct{} {
	var accountsMap = make(map[string]struct{}, len(accounts))
	for _, account := range accounts {
		accountsMap[account.Address] = struct{}{}
	}
	return accountsMap
}
