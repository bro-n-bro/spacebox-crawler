package bank

import (
	"context"
	"encoding/json"
	"fmt"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func (m *Module) HandleGenesis(_ context.Context, _ *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {

	// Read the genesis state
	var genState govtypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[govtypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading gov genesis data: %s", err)
	}

	// TODO:
	_ = genState.Proposals

	return nil
}
