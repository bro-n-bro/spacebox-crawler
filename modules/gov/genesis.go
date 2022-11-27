package bank

import (
	"context"
	"encoding/json"
	"fmt"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	tmtypes "github.com/tendermint/tendermint/types"
)

func (m *Module) HandleGenesis(_ context.Context, _ *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {

	// Read the genesis state
	var genState govtypesv1beta1.GenesisState
	err := m.cdc.UnmarshalJSON(appState[govtypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading gov genesis data: %s", err)
	}

	// TODO:
	_ = genState.Proposals

	return nil
}
