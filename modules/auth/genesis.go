package auth

import (
	"context"
	"encoding/json"

	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/hexy-dev/spacebox-crawler/modules/utils"
	"github.com/hexy-dev/spacebox/broker/model"
)

func (m *Module) HandleGenesis(ctx context.Context, _ *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	accounts, err := utils.GetGenesisAccounts(appState, m.cdc)
	if err != nil {
		return err
	}

	for _, acc := range accounts {
		// TODO: test it
		if err = m.broker.PublishAccount(ctx, model.Account{
			Address: acc.Address,
			Height:  acc.Height,
		}); err != nil {
			return err
		}
	}

	return nil
}
