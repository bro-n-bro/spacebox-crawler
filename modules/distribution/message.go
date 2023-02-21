package distribution

import (
	"context"
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/pkg/errors"

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

		coin, err := coinFromBytes([]byte(value))
		if err != nil {
			return fmt.Errorf("%w failed to convert %s string to coin", err, value)
		}

		// TODO: test it
		return m.broker.PublishDelegationRewardMessage(ctx, model.DelegationRewardMessage{
			Coin:             m.tbM.MapCoin(coin),
			Height:           tx.Height,
			DelegatorAddress: msg.DelegatorAddress,
			ValidatorAddress: msg.ValidatorAddress,
			TxHash:           tx.TxHash,
			MsgIndex:         int64(index),
		})
	}

	return nil
}

// coinFromBytes converts slice bytes to coin type
// loop over value with index and find first not digit byte. it means that before the index - digits and after - chars
// ex: 123abc - > amount: 123, denom: abc
func coinFromBytes(value []byte) (types.Coin, error) {
	for i := 0; i < len(value); i++ {
		v := value[i]
		if v < 48 || v > 57 { // bytes range [48-57] == digits range [0-9]
			amount, err := strconv.ParseFloat(string(value[:i]), 64)
			if err != nil {
				return types.Coin{}, err
			}
			return types.Coin{
				Denom:  string(value[i:]),
				Amount: amount,
			}, nil
		}
	}

	return types.Coin{}, errors.New("not found")
}
