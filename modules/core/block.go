package core

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/hexy-dev/spacebox-crawler/types"
)

func (m *Module) HandleBlock(ctx context.Context, block *types.Block, vals *tmctypes.ResultValidators) error {
	if err := m.broker.PublishBlock(ctx, m.tbM.MapBlock(block, block.TotalGas)); err != nil {
		m.log.Error().
			Err(err).
			Int64("block_height", block.Height).
			Msg("PublishBlock error")
		return err
	}
	return m.publishValidators(ctx, vals.Validators)
}

func (m *Module) publishValidators(ctx context.Context, vals []*tmtypes.Validator) error {
	var validators = make([]*types.Validator, len(vals))
	for index, val := range vals {
		consAddr := sdk.ConsAddress(val.Address).String()

		consPubKey, err := types.ConvertValidatorPubKeyToBech32String(val.PubKey)
		if err != nil {
			return fmt.Errorf("failed to convert validator public key for validators %s: %s", consAddr, err)
		}

		validators[index] = types.NewValidator(consAddr, consPubKey)
	}

	// TODO: save to mongo?
	// TODO: save it?
	return m.broker.PublishValidators(ctx, m.tbM.MapValidators(validators))

}
