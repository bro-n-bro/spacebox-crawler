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
	if err := m.parseProposerRewardEvent(ctx, eventsMap, height); err != nil {
		m.log.Error().
			Err(err).
			Str("handler", "HandleBeginBlocker").
			Int64("height", height).
			Msg("failed to parse proposer reward event")
		return err
	}

	if err := m.parseCommissionEvent(ctx, eventsMap, height); err != nil {
		m.log.Error().
			Err(err).
			Str("handler", "HandleBeginBlocker").
			Int64("height", height).
			Msg("failed to parse distribution commission event")
		return err
	}

	if err := m.parseRewardsEvent(ctx, eventsMap, height); err != nil {
		m.log.Error().
			Err(err).
			Str("handler", "HandleBeginBlocker").
			Int64("height", height).
			Msg("failed to parse distribution rewards event")
		return err
	}

	return nil
}

// parseProposerRewardEvent parses proposer reward event.
func (m *Module) parseProposerRewardEvent(ctx context.Context, eventsMap types.BlockerEvents, height int64) error {
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
			m.log.Warn().Str("func", "parseProposerRewardEvent").Msg("not enough attributes in event")
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
						Str("func", "parseProposerRewardEvent").
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
				Str("func", "parseProposerRewardEvent").
				Int64("height", height).
				Msg("error while publishing proposer reward")
			return err
		}
	}

	return nil
}

// parseCommissionEvent parses distribution commission event.
// nolint: dupl
func (m *Module) parseCommissionEvent(ctx context.Context, eventsMap types.BlockerEvents, height int64) error {
	events, ok := eventsMap[distrtypes.EventTypeCommission]
	if !ok {
		return nil
	}

	var (
		validator string
		coin      model.Coin
	)

	for _, event := range events {
		if len(event.Attributes) < 2 {
			m.log.Warn().Str("func", "parseCommissionEvent").Msg("not enough attributes in event")
			continue
		}

		for _, attr := range event.Attributes {
			switch string(attr.Key) {
			case sdk.AttributeKeyAmount:
				coins, err := utils.ParseCoinsFromString(string(attr.Value))
				if err != nil {
					m.log.Error().
						Err(err).
						Str("func", "parseProposerRewardEvent").
						Int64("height", height).
						Msg("failed to convert string to coins by commissionEvent")

					return fmt.Errorf("failed to convert %q to coin: %w", string(attr.Value), err)
				}

				if len(coins) > 0 {
					coin = m.tbM.MapCoin(coins[0])
				}
			case distrtypes.AttributeKeyValidator:
				validator = string(attr.Value)
			}
		}

		if err := m.broker.PublishDistributionCommission(ctx, model.DistributionCommission{
			Height:    height,
			Validator: validator,
			Amount:    coin,
		}); err != nil {
			m.log.Error().
				Err(err).
				Str("func", "parseCommissionEvent").
				Int64("height", height).
				Msg("error while publishing distribution commission")
			return err
		}
	}

	return nil
}

// parseRewardsEvent parses rewards event.
// nolint:dupl
func (m *Module) parseRewardsEvent(ctx context.Context, eventsMap types.BlockerEvents, height int64) error {
	events, ok := eventsMap[distrtypes.EventTypeRewards]
	if !ok {
		return nil
	}

	var (
		validator string
		coin      model.Coin
	)

	for _, event := range events {
		if len(event.Attributes) < 2 {
			m.log.Warn().Str("func", "parseRewardsEvent").Msg("not enough attributes in event")
			continue
		}

		for _, attr := range event.Attributes {
			switch string(attr.Key) {
			case sdk.AttributeKeyAmount:
				coins, err := utils.ParseCoinsFromString(string(attr.Value))
				if err != nil {
					m.log.Error().
						Err(err).
						Str("func", "parseRewardsEvent").
						Int64("height", height).
						Msg("failed to convert string to coins by rewardsEvent")

					return fmt.Errorf("failed to convert %q to coin: %w", string(attr.Value), err)
				}

				if len(coins) > 0 {
					coin = m.tbM.MapCoin(coins[0])
				}
			case distrtypes.AttributeKeyValidator:
				validator = string(attr.Value)
			}
		}

		if err := m.broker.PublishDistributionReward(ctx, model.DistributionReward{
			Height:    height,
			Validator: validator,
			Amount:    coin,
		}); err != nil {
			m.log.Error().
				Err(err).
				Str("func", "parseRewardsEvent").
				Int64("height", height).
				Msg("error while publishing distribution reward")
			return err
		}
	}

	return nil
}
