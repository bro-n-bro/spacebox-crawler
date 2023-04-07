package distribution

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/bro-n-bro/spacebox-crawler/modules/utils"
	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

func (m *Module) HandleBeginBlocker(ctx context.Context, eventsMap types.BlockerEvents, height int64) error {
	events, ok := eventsMap[distrtypes.EventTypeProposerReward]
	if !ok {
		return nil
	}

	var (
		coin      model.Coin
		validator string
	)

	for _, event := range events {
		if len(event.Attributes) < 2 {
			m.log.Warn().Str("handler", "HandleBeginBlocker").Msg("not enough attributes in event")
			continue
		}

		for _, attr := range event.Attributes {
			switch string(attr.Key) {
			case distrtypes.AttributeKeyValidator:
				validator = string(attr.Value)
			case sdk.AttributeKeyAmount:
				coins, err := utils.ParseCoinsFromString(string(attr.Value))
				if err != nil {
					m.log.Error().
						Err(err).
						Str("handler", "HandleBeginBlocker").
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
			m.log.Error().
				Err(err).
				Str("handler", "HandleBeginBlocker").
				Int64("height", height).
				Msg("error while publishing proposer reward")
			return err
		}
	}

	return nil
}
