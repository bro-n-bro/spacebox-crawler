package distribution

import (
	"context"
	"encoding/base64"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/bro-n-bro/spacebox-crawler/modules/utils"
	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

var (
	base64KeyValidator = base64.StdEncoding.EncodeToString([]byte(distrtypes.AttributeKeyValidator))
	base64KeyAmount    = base64.StdEncoding.EncodeToString([]byte(sdk.AttributeKeyAmount))
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
		err       error
	)

	for _, event := range events {
		if len(event.Attributes) < 2 {
			m.log.Warn().Str("func", "parseProposerRewardEvent").Msg("not enough attributes in event")
			continue
		}

		validator, coin, err = m.parseAttributes(event)
		if err != nil {
			m.log.Error().
				Err(err).
				Str("func", "parseProposerRewardEvent").
				Int64("height", height).
				Msg("failed to parse attributes")
			return err
		}

		if err = m.broker.PublishProposerReward(ctx, model.ProposerReward{
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
func (m *Module) parseCommissionEvent(ctx context.Context, eventsMap types.BlockerEvents, height int64) error {
	events, ok := eventsMap[distrtypes.EventTypeCommission]
	if !ok {
		return nil
	}

	var (
		validator string
		coin      model.Coin
		err       error
	)

	for _, event := range events {
		if len(event.Attributes) < 2 {
			m.log.Warn().
				Str("event", distrtypes.EventTypeCommission).
				Int64("height", height).
				Msg("not enough attributes in event")
			continue
		}

		validator, coin, err = m.parseAttributes(event)
		if err != nil {
			m.log.Error().
				Err(err).
				Str("func", "parseCommissionEvent").
				Int64("height", height).
				Msg("failed to parse attributes")
			return err
		}

		if err = m.broker.PublishDistributionCommission(ctx, model.DistributionCommission{
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
func (m *Module) parseRewardsEvent(ctx context.Context, eventsMap types.BlockerEvents, height int64) error {
	events, ok := eventsMap[distrtypes.EventTypeRewards]
	if !ok {
		return nil
	}

	var (
		validator string
		coin      model.Coin
		err       error
	)

	for _, event := range events {
		if len(event.Attributes) < 2 {
			m.log.Warn().Str("func", "parseRewardsEvent").Msg("not enough attributes in event")
			continue
		}

		validator, coin, err = m.parseAttributes(event)
		if err != nil {
			m.log.Error().
				Err(err).
				Str("func", "parseCommissionEvent").
				Int64("height", height).
				Msg("failed to parse attributes")
			return err
		}

		if err = m.broker.PublishDistributionReward(ctx, model.DistributionReward{
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

func (m *Module) parseAttributes(event abci.Event) (string, model.Coin, error) {
	var (
		coin      model.Coin
		validator string
	)

	for _, attr := range event.Attributes {
		// try to decode value if needed
		switch attr.Key {
		case base64KeyValidator, base64KeyAmount:
			var err error
			attr.Value, err = utils.DecodeToString(attr.Value)
			if err != nil {
				return "", model.Coin{}, err
			}
		}

		switch attr.Key {
		case distrtypes.AttributeKeyValidator, base64KeyValidator:
			validator = attr.Value
		case sdk.AttributeKeyAmount, base64KeyAmount:
			coins, err := utils.ParseCoinsFromString(attr.Value)
			if err != nil {
				return "", model.Coin{}, fmt.Errorf("failed to convert %q to coin: %w", attr.Value, err)
			}

			if len(coins) > 0 {
				coin = m.tbM.MapCoin(coins[0])
			}
		}
	}

	return validator, coin, nil
}
