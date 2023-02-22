package distribution

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

var (
	withdrawDelegationRewardRegex = regexp.MustCompile(`^(\-?[0-9]+(\.[0-9]+)?)([0-9a-zA-Z/]+)$`)
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
			coin, err = coinsFromAttribute(value)
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
	}

	return nil
}

// coinsFromAttribute converts string to coin type
func coinsFromAttribute(value string) (types.Coins, error) {
	rows := strings.Split(value, ",")
	res := make(types.Coins, len(rows))
	for i, row := range rows {
		bits := withdrawDelegationRewardRegex.FindStringSubmatch(row)
		if len(bits) < 4 {
			continue
		}
		amount, err := strconv.ParseFloat(bits[1], 64)
		if err != nil {
			return types.Coins{}, errors.Wrap(err, "failed to parse float")
		}
		res[i] = types.Coin{
			Denom:  bits[3],
			Amount: amount,
		}
	}
	return res, nil
}
