package slashing

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"

	"github.com/bro-n-bro/spacebox-crawler/modules/utils"
	tb "github.com/bro-n-bro/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

var (
	errCantFindBurnedCoin = errors.New("can't find burned tokens")

	base64KeyAddress     = base64.StdEncoding.EncodeToString([]byte(slashingtypes.AttributeKeyAddress))
	base64KeyPower       = base64.StdEncoding.EncodeToString([]byte(slashingtypes.AttributeKeyPower))
	base64KeyReason      = base64.StdEncoding.EncodeToString([]byte(slashingtypes.AttributeKeyReason))
	base64KeyJailed      = base64.StdEncoding.EncodeToString([]byte(slashingtypes.AttributeKeyJailed))
	base64KeyBurnedCoins = base64.StdEncoding.EncodeToString([]byte(slashingtypes.AttributeKeyBurnedCoins))
	base64KeyAmount      = base64.StdEncoding.EncodeToString([]byte(sdk.AttributeKeyAmount))
)

func (m *Module) HandleBeginBlocker(ctx context.Context, eventsMap types.BlockerEvents, height int64) error {
	return m.handleSlashEvent(ctx, eventsMap, height)
}

// nolint:gocognit
func (m *Module) handleSlashEvent(ctx context.Context, eventsMap types.BlockerEvents, height int64) error {
	events, ok := eventsMap[slashingtypes.EventTypeSlash]
	if !ok {
		return nil
	}

	var operatorAddress, power, reason, jailed string
	for _, e := range events {
		if len(e.Attributes) < 4 {
			m.log.Warn().
				Str("event", slashingtypes.EventTypeSlash).
				Int64("height", height).
				Msg("not enough attributes in event")
			continue
		}

		var burnedCoin *model.Coin
		for _, attr := range e.Attributes {
			// try to decode value if needed
			switch attr.Key {
			case base64KeyAddress, base64KeyPower, base64KeyReason, base64KeyJailed, base64KeyBurnedCoins:
				var err error
				attr.Value, err = utils.DecodeToString(attr.Value)
				if err != nil {
					return err
				}
			}

			switch attr.Key {
			case slashingtypes.AttributeKeyAddress, base64KeyAddress: // required
				operatorAddress = attr.Value
			case slashingtypes.AttributeKeyPower, base64KeyPower: // required
				power = attr.Value
			case slashingtypes.AttributeKeyReason, base64KeyReason: // required
				reason = attr.Value
			case slashingtypes.AttributeKeyJailed, base64KeyJailed: // required
				jailed = attr.Value
			case slashingtypes.AttributeKeyBurnedCoins, base64KeyBurnedCoins: // not required
				coins, err := utils.ParseCoinsFromString(attr.Value)
				if err != nil {
					m.log.Error().
						Err(err).
						Str("func", "parseProposerRewardEvent").
						Int64("height", height).
						Msg("failed to convert string to coins by commissionEvent")

					return fmt.Errorf("failed to convert %q to coin: %w", attr.Value, err)
				}
				if len(coins) > 0 {
					coin := m.tbM.MapCoin(coins[0])
					burnedCoin = &coin
				}
			}
		}

		var burned model.Coin
		if burnedCoin == nil {
			var err error
			burned, err = getCoin(eventsMap, m.tbM)
			if err != nil {
				if errors.Is(err, errCantFindBurnedCoin) {
					m.log.Warn().
						Str("event", banktypes.EventTypeCoinBurn).
						Int64("height", height).
						Msg(err.Error())
					continue
				}

				m.log.Error().
					Str("event", banktypes.EventTypeCoinBurn).
					Int64("height", height).
					Msg(err.Error())
				return err
			}
		} else {
			burned = *burnedCoin
		}

		if err := m.broker.PublishHandleValidatorSignature(ctx, model.HandleValidatorSignature{
			OperatorAddress: operatorAddress,
			Power:           power,
			Reason:          reason,
			Jailed:          jailed,
			Burned:          burned,
			Height:          height,
		}); err != nil {
			return err
		}
	}

	return nil
}

func getCoin(eventsMap types.BlockerEvents, mapper tb.ToBroker) (model.Coin, error) {
	bankEvents, ok := eventsMap[banktypes.EventTypeCoinBurn]
	if !ok || len(bankEvents) == 0 || len(bankEvents[0].Attributes) < 2 {
		// burned tokens not found in any events
		return model.Coin{}, errCantFindBurnedCoin
	}

	var res model.Coin
	for _, bankAttr := range bankEvents[0].Attributes {
		if bankAttr.Key != sdk.AttributeKeyAmount && bankAttr.Key != base64KeyAmount {
			continue
		}

		if bankAttr.Key == base64KeyAmount {
			var err error
			bankAttr.Value, err = utils.DecodeToString(bankAttr.Value)
			if err != nil {
				return model.Coin{}, err
			}
		}

		coins, err := utils.ParseCoinsFromString(bankAttr.Value)
		if err != nil {
			err = fmt.Errorf("failed to convert %q to coin: %w", bankAttr.Value, err)
			return model.Coin{}, err
		}

		if len(coins) > 0 {
			res = mapper.MapCoin(coins[0])
			break
		}
	}

	return res, nil
}
