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
func (m *Module) HandleGenesis(ctx context.Context, doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	var bankState banktypes.GenesisState
	if err := m.cdc.UnmarshalJSON(appState[banktypes.ModuleName], &bankState); err != nil {
		return err
	}

	// Store the balances
	accounts, err := utils.GetGenesisAccounts(appState, m.cdc)
	if err != nil {
		return err
	}

	uniqueAddresses := make(map[string]struct{})
	for _, acc := range accounts {
		uniqueAddresses[acc.Address] = struct{}{}
	}

	for _, balance := range bankState.Balances {
		_, ok := uniqueAddresses[balance.Address]
		if !ok {
			// skip already published
			continue
		}
		ab := types.NewAccountBalance(balance.Address, balance.Coins, doc.InitialHeight)
		// TODO: test it
		if err = m.broker.PublishAccountBalance(ctx, m.tbM.MapAccountBalance(ab)); err != nil {
			return err
		}
	}

	return nil
}
