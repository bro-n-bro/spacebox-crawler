package distribution

import (
	"context"
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/bro-n-bro/spacebox-crawler/modules/utils"
	"github.com/bro-n-bro/spacebox/broker/model"
)

const (
	proposalRewardEvent = "proposer_reward"
)

func (m *Module) HandleBeginBlocker(ctx context.Context, events []abci.Event, height int64) error {
	for _, e := range events {
		if e.Type == proposalRewardEvent {
			if len(e.Attributes) < 2 {
				m.log.Warn().Int64("height", height).Msg("proposer_reward event less than 2 attributes")

				// because it contains only one event
				break
			}

			var (
				coin      model.Coin
				validator string
			)

			for _, attr := range e.Attributes {
				switch string(attr.Key) {
				case "validator":
					validator = string(attr.Value)
				case "amount":
					coins, err := utils.ParseCoinsFromString(string(attr.Value))
					if err != nil {
						m.log.Error().
							Err(err).
							Int64("height", height).
							Msg("failed to convert string to coins by proposalRewardEvent")

						return fmt.Errorf("failed to convert %q to coin: %w", string(attr.Value), err)
					}

					if len(coins) > 0 {
						coin = m.tbM.MapCoin(coins[0])
					}
				}
			}

			if err := m.broker.PublishProposerReward(ctx, model.ProposerReward{
				Height:    height,
				Validator: validator,
				Reward:    coin,
			}); err != nil {
				return err
			}

			// because it contains only one event
			break
		}
	}
	return nil
}
