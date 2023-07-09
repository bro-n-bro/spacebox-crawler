package auth

import (
	"context"
	"encoding/json"
	"strings"

	cometbfttypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bro-n-bro/spacebox-crawler/modules/utils"
	"github.com/bro-n-bro/spacebox/broker/model"
)

func (m *Module) HandleGenesis(ctx context.Context, _ *cometbfttypes.GenesisDoc, appState map[string]json.RawMessage) error {
	accounts, err := utils.GetGenesisAccounts(appState, m.cdc)
	if err != nil {
		return err
	}

	prefix := sdk.GetConfig().GetBech32AccountAddrPrefix()
	for _, acc := range accounts {
		if strings.HasPrefix(acc.Address, prefix) {
			// TODO: test it
			if err = m.broker.PublishAccount(ctx, model.Account{
				Address: acc.Address,
				Height:  acc.Height,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}
