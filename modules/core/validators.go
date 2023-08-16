package core

import (
	"context"

	cometbftcoretypes "github.com/cometbft/cometbft/rpc/core/types"

	"github.com/bro-n-bro/spacebox/broker/model"
)

// HandleValidators handles validators for each block height.
func (m *Module) HandleValidators(ctx context.Context, vals *cometbftcoretypes.ResultValidators) error {
	for _, val := range vals.Validators {
		if err := m.broker.PublishValidatorVotingPower(ctx, model.ValidatorVotingPower{
			Height:           vals.BlockHeight,
			VotingPower:      val.VotingPower,
			ValidatorAddress: string(val.Address),
			ProposerPriority: val.ProposerPriority,
		}); err != nil {
			return err
		}
	}

	return nil
}
