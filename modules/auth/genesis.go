package auth

import (
	"context"
	"encoding/json"

	cometbfttypes "github.com/cometbft/cometbft/types"

	"github.com/bro-n-bro/spacebox-crawler/modules/utils"
)

func (m *Module) HandleGenesis(
	ctx context.Context,
	doc *cometbfttypes.GenesisDoc,
	appState map[string]json.RawMessage,
) error {

	accounts, err := utils.GetGenesisAccounts(appState, m.cdc)
	if err != nil {
		return err
	}

	for _, acc := range accounts {
		if err = utils.GetAndPublishAccount(ctx, acc.Address, doc.InitialHeight, m.accCache, m.broker,
			m.client.AuthQueryClient); err != nil {
			return err
		}
	}

	return nil
}
