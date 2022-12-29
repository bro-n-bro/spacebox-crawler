package bank

import (
	"context"
	"encoding/json"
	"fmt"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	tmtypes "github.com/tendermint/tendermint/types"

	govutils "github.com/hexy-dev/spacebox-crawler/modules/gov/utils"
)

func (m *Module) HandleGenesis(ctx context.Context, _ *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {

	// Read the genesis state
	var genState govtypesv1beta1.GenesisState
	err := m.cdc.UnmarshalJSON(appState[govtypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading gov genesis data: %s", err)
	}

	proposals := genState.Proposals
	if err = govutils.SaveProposals(ctx, proposals, m.broker, m.tbM); err != nil {
		return fmt.Errorf("error while saving genesis proposal data: %s", err)
	}

	return nil
}
