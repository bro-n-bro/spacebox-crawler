package mint

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

func (m *Module) HandleBeginBlocker(ctx context.Context, eventsMap types.BlockerEvents, height int64) error {
	events, ok := eventsMap[minttypes.EventTypeMint]
	if !ok {
		return nil
	}

	var (
		bondedRatio, inflation, annualProvisions float64
		amount                                   int64
		err                                      error
	)

	for _, event := range events {
		if len(event.Attributes) < 4 {
			m.log.Warn().Str("handler", "HandleBeginBlocker").Msg("not enough attributes in event")
			continue
		}

		for _, attr := range event.Attributes {
			switch string(attr.Key) {
			case minttypes.AttributeKeyBondedRatio:
				bondedRatio, err = strconv.ParseFloat(string(attr.Value), 64)
				if err != nil {
					m.log.Error().
						Err(err).
						Str("handler", "HandleBeginBlocker").
						Int64("height", height).
						Msg("failed to convert string to float64 for AttributeKeyBondedRatio")
					return err
				}
			case minttypes.AttributeKeyInflation:
				inflation, err = strconv.ParseFloat(string(attr.Value), 64)
				if err != nil {
					m.log.Error().
						Err(err).
						Str("handler", "HandleBeginBlocker").
						Int64("height", height).
						Msg("failed to convert string to float64 for AttributeKeyBondedRatio")
					return err
				}
			case minttypes.AttributeKeyAnnualProvisions:
				annualProvisions, err = strconv.ParseFloat(string(attr.Value), 64)
				if err != nil {
					m.log.Error().
						Err(err).
						Str("handler", "HandleBeginBlocker").
						Int64("height", height).
						Msg("failed to convert string to float64 for AttributeKeyBondedRatio")
					return err
				}
			case sdk.AttributeKeyAmount:
				amount, err = strconv.ParseInt(string(attr.Value), 10, 64)
				if err != nil {
					m.log.Error().
						Err(err).
						Str("handler", "HandleBeginBlocker").
						Int64("height", height).
						Msg("failed to convert string to int64 for AttributeKeyAmount")
					return err
				}
			}
		}

		if err = m.broker.PublishAnnualProvision(ctx, model.AnnualProvision{
			Height:          height,
			Amount:          amount,
			AnnualProvision: annualProvisions,
			BondedRatio:     bondedRatio,
			Inflation:       inflation,
		}); err != nil {
			m.log.Error().
				Err(err).
				Str("handler", "HandleBeginBlocker").
				Int64("height", height).
				Msg("failed to publish annual provision")
			continue
		}
	}

	return nil
}
