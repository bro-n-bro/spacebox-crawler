package bank

import (
	"context"
	"encoding/json"

	cometbfttypes "github.com/cometbft/cometbft/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/bro-n-bro/spacebox-crawler/modules/utils"
	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// HandleGenesis handles the genesis state of the x/bank module in order to store the initial values
// of the different account balances.
func (m *Module) HandleGenesis(
	ctx context.Context,
	doc *cometbfttypes.GenesisDoc,
	appState map[string]json.RawMessage) error {

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
		if _, ok := uniqueAddresses[balance.Address]; !ok { // skip already published
			continue
		}

		// TODO: test it
		if err = m.broker.PublishAccountBalance(ctx, model.AccountBalance{
			Address: balance.Address,
			Height:  doc.InitialHeight,
			Coins:   m.tbM.MapCoins(types.NewCoinsFromCdk(balance.Coins)),
		}); err != nil {
			return err
		}
	}

	return nil
}
