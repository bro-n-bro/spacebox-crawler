package auth

import (
	"context"
	"encoding/json"

	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/hexy-dev/spacebox-crawler/modules/utils"
)

func (m *Module) HandleGenesis(ctx context.Context, _ *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	accounts, err := utils.GetGenesisAccounts(appState, m.cdc)
	if err != nil {
		return err
	}

	// TODO: test it
	if err = m.broker.PublishAccounts(ctx, m.tbM.MapAccounts(accounts)); err != nil {
		return err
	}
	return nil
}
