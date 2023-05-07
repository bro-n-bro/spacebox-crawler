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

func (m *Module) HandleMessage(ctx context.Context, index int, cosmosMsg sdk.Msg, tx *types.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch msg := cosmosMsg.(type) {
	case *distrtypes.MsgWithdrawDelegatorReward:
		event, err := tx.FindEventByType(index, distrtypes.EventTypeWithdrawRewards)
		if err != nil {
			return err
		}

		value, err := tx.FindAttributeByKey(event, "amount")
		if err != nil {
			return err
		}

		coin := types.Coins{}
		if value != "" {
			coin, err = utils.ParseCoinsFromString(value)
			if err != nil {
				m.log.Error().
					Err(err).
					Str("tx_hash", tx.TxHash).
					Msg("failed to convert string to coin by MsgWithdrawDelegatorReward")
				return fmt.Errorf("%w failed to convert %s string to coin", err, value)
			}
		}

		// TODO: test it
		return m.broker.PublishDelegationRewardMessage(ctx, model.DelegationRewardMessage{
			Coins:            m.tbM.MapCoins(coin),
			Height:           tx.Height,
			DelegatorAddress: msg.DelegatorAddress,
			ValidatorAddress: msg.ValidatorAddress,
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
			ValidatorAddress:   msg.ValidatorAddress,
		})
	}

	return nil
}
