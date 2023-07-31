package mint

import (
	"context"
	"encoding/base64"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"github.com/bro-n-bro/spacebox-crawler/modules/utils"
	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

var (
	base64KeyBondedRatio      = base64.StdEncoding.EncodeToString([]byte(minttypes.AttributeKeyBondedRatio))
	base64KeyInflation        = base64.StdEncoding.EncodeToString([]byte(minttypes.AttributeKeyInflation))
	base64KeyAnnualProvisions = base64.StdEncoding.EncodeToString([]byte(minttypes.AttributeKeyAnnualProvisions))
	base64KeyAmount           = base64.StdEncoding.EncodeToString([]byte(sdk.AttributeKeyAmount))
)

//nolint:gocognit
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
			// try to decode value if needed
			switch attr.Key {
			case base64KeyBondedRatio, base64KeyInflation, base64KeyAnnualProvisions, base64KeyAmount:
				attr.Value, err = utils.DecodeToString(attr.Value)
				if err != nil {
					return err
				}
			}

			switch attr.Key {
			case minttypes.AttributeKeyBondedRatio, base64KeyBondedRatio:
				bondedRatio, err = strconv.ParseFloat(attr.Value, 64)
				if err != nil {
					m.log.Error().
						Err(err).
						Str("handler", "HandleBeginBlocker").
						Int64("height", height).
						Msg("failed to convert string to float64 for AttributeKeyBondedRatio")
					return err
				}
			case minttypes.AttributeKeyInflation, base64KeyInflation:
				inflation, err = strconv.ParseFloat(attr.Value, 64)
				if err != nil {
					m.log.Error().
						Err(err).
						Str("handler", "HandleBeginBlocker").
						Int64("height", height).
						Msg("failed to convert string to float64 for AttributeKeyBondedRatio")
					return err
				}
			case minttypes.AttributeKeyAnnualProvisions, base64KeyAnnualProvisions:
				annualProvisions, err = strconv.ParseFloat(attr.Value, 64)
				if err != nil {
					m.log.Error().
						Err(err).
						Str("handler", "HandleBeginBlocker").
						Int64("height", height).
						Msg("failed to convert string to float64 for AttributeKeyBondedRatio")
					return err
				}
			case sdk.AttributeKeyAmount, base64KeyAmount:
				amount, err = strconv.ParseInt(attr.Value, 10, 64)
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
			Height:           height,
			Amount:           amount,
			AnnualProvisions: annualProvisions,
			BondedRatio:      bondedRatio,
			Inflation:        inflation,
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
