package core

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/hexy-dev/spacebox-crawler/types"
	"github.com/hexy-dev/spacebox/broker/model"
)

func (m *Module) ValidatorsHandler(ctx context.Context, vals *tmctypes.ResultValidators) error {
	for _, val := range vals.Validators {
		consAddr := sdk.ConsAddress(val.Address).String()

		consPubKey, err := types.ConvertValidatorPubKeyToBech32String(val.PubKey)
		if err != nil {
			msg := fmt.Sprintf("failed to convert validator public key for validators %s: %s", consAddr, err)
			return errors.New(msg)
		}

		// TODO: save to mongo?
		// TODO: save it?
		if err = m.broker.PublishValidator(ctx, model.Validator{
			ConsensusAddress: consAddr,
			ConsensusPubkey:  consPubKey,
		}); err != nil {
			return err
		}
	}

	return nil
}
