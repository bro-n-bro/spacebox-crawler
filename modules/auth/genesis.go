package auth

import (
	"context"
	"encoding/json"

	tmtypes "github.com/tendermint/tendermint/types"

	"bro-n-bro-osmosis/modules/utils"
)

func (m *Module) HandleGenesis(ctx context.Context, doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	accounts, err := utils.GetGenesisAccounts(appState, m.cdc)
	if err != nil {
		return err
	}

	// TODO:
	_ = accounts
	return nil
}
