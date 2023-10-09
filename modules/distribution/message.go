package distribution

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/bro-n-bro/spacebox-crawler/modules/utils"
	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

var (
	errNotFoundTxResult = errors.New("not found tx result by index")
)

func (m *Module) HandleMessage(ctx context.Context, index int, cosmosMsg sdk.Msg, tx *types.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch msg := cosmosMsg.(type) {
	case *distrtypes.MsgWithdrawDelegatorReward:
		coin, err := m.findCoinFromEventByValidator(ctx, tx, index, msg.ValidatorAddress)
		if err != nil {
			if !errors.Is(err, types.ErrNoEventFound) {
				return err
			}
			// pass zero coins
		}

		return m.broker.PublishDelegationRewardMessage(ctx, model.DelegationRewardMessage{
			Coins:            m.tbM.MapCoins(coin),
			Height:           tx.Height,
			DelegatorAddress: msg.DelegatorAddress,
			OperatorAddress:  msg.ValidatorAddress,
			TxHash:           tx.TxHash,
			MsgIndex:         int64(index),
		})
	case *distrtypes.MsgSetWithdrawAddress:
		return m.broker.PublishSetWithdrawAddressMessage(ctx, model.SetWithdrawAddressMessage{
			Height:           tx.Height,
			DelegatorAddress: msg.DelegatorAddress,
			WithdrawAddress:  msg.WithdrawAddress,
			TxHash:           tx.TxHash,
			MsgIndex:         int64(index),
		})
	case *distrtypes.MsgWithdrawValidatorCommission:
		event, err := tx.FindEventByType(index, distrtypes.EventTypeWithdrawCommission)
		if err != nil {
			return err
		}

		value, err := tx.FindAttributeByKey(event, sdk.AttributeKeyAmount)
		if err != nil {
			return err
		}

		coins := types.Coins{}
		if value != "" {
			coins, err = utils.ParseCoinsFromString(value)
			if err != nil {
				m.log.Error().
					Err(err).
					Str("tx_hash", tx.TxHash).
					Int64("height", tx.Height).
					Msg("failed to convert string to coin by MsgWithdrawValidatorCommission height")
				return fmt.Errorf("%w failed to convert %s string to coin height:%v", err, value, tx.Height)
			}
		}
		return m.broker.PublishWithdrawValidatorCommissionMessage(ctx, model.WithdrawValidatorCommissionMessage{
			Height:             tx.Height,
			TxHash:             tx.TxHash,
			MsgIndex:           int64(index),
			WithdrawCommission: m.tbM.MapCoins(coins),
			OperatorAddress:    msg.ValidatorAddress,
		})
	}

	return nil
}

func (m *Module) findCoinFromEventByValidator(
	ctx context.Context,
	tx *types.Tx,
	index int,
	validatorAddress string,
) (types.Coins, error) {

	found := false
	coins := types.Coins{}

Events:
	for _, ev := range tx.Logs[index].Events {
		if ev.Type != distrtypes.EventTypeWithdrawRewards {
			continue
		}

		for _, attr := range ev.Attributes {
			if attr.Key == distrtypes.AttributeKeyValidator &&
				compareValueInBase64(validatorAddress, attr.Value) {

				found = true
			}

			if found && attr.Key == sdk.AttributeKeyAmount {
				var err error
				coins, err = utils.ParseCoinsFromString(attr.Value)
				if err != nil {
					return coins, fmt.Errorf("%w failed to convert %s string to coin", err, attr.Value)
				}

				break Events
			}
		}
	}

	if !found {
		var err error
		coins, err = m.findWithdrawRewardInTxEvents(ctx, validatorAddress, index, tx)
		if err != nil {
			m.log.Warn().
				Str("tx_hash", tx.TxHash).
				Int64("height", tx.Height).
				Str("event", distrtypes.EventTypeWithdrawRewards).
				Msg("not found event in tx")

			return coins, fmt.Errorf("%w: %s inside tx with hash %s", types.ErrNoEventFound,
				distrtypes.EventTypeWithdrawRewards, tx.TxHash)
		}
	}

	return coins, nil
}

func (m *Module) findWithdrawRewardInTxEvents(
	ctx context.Context,
	address string,
	msgIndex int,
	tx *types.Tx,
) (types.Coins, error) {

	txResults, err := m.rpcCli.GetTxResults(ctx, tx.Height)
	if err != nil {
		return nil, err
	}

	if len(txResults)-1 < msgIndex {
		return nil, fmt.Errorf("%w: %d", errNotFoundTxResult, msgIndex)
	}

	for _, event := range txResults[msgIndex].Events {
		if event.Type != distrtypes.EventTypeWithdrawRewards {
			continue
		}

		for _, attr := range event.Attributes {
			if (attr.Key == distrtypes.AttributeKeyValidator || attr.Key == base64KeyValidator) &&
				compareValueInBase64(address, attr.Value) {

				for _, attr2 := range event.Attributes {
					if attr2.Key == sdk.AttributeKeyAmount || attr2.Key == base64KeyAmount {
						coins, err := utils.ParseCoinsFromString(attr.Value)
						if err != nil {
							//nolint:lll
							m.log.Error().
								Err(err).
								Str("tx_hash", tx.TxHash).
								Msg("failed to convert string to coin by MsgWithdrawDelegatorReward")

							//nolint:lll
							return coins, fmt.Errorf("%w failed to convert %s string to coin", err, attr.Value)
						}

						return coins, nil
					}
				}
			}
		}
	}

	return nil, types.ErrNoEventFound
}

func compareValueInBase64(source, target string) bool {
	if source == target {
		return true
	}

	val, err := base64.StdEncoding.DecodeString(target)
	if err != nil {
		return false
	}

	return source == string(val)
}
