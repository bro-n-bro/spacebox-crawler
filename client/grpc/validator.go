package grpc

import (
	"context"

	"github.com/cometbft/cometbft/crypto/ed25519"
	cometbftcoretypes "github.com/cometbft/cometbft/rpc/core/types"
	cometbfttypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/types/query"
)

const (
	defaultLimit = 150
)

func (c *Client) Validators(ctx context.Context, height int64) (*cometbftcoretypes.ResultValidators, error) {
	vals := &cometbftcoretypes.ResultValidators{
		BlockHeight: height,
	}

	var (
		nextKey []byte
	)

	for {
		respPb, err := c.TmsService.GetValidatorSetByHeight(ctx, &tmservice.GetValidatorSetByHeightRequest{
			Height: height,
			Pagination: &query.PageRequest{
				Key:        nextKey,
				Limit:      defaultLimit,
				CountTotal: true,
			},
		})
		if err != nil {
			return nil, err
		}

		nextKey = respPb.Pagination.NextKey

		vals.Total = int(respPb.Pagination.Total)
		if len(nextKey) == 0 { // first iteration
			vals.Validators = make([]*cometbfttypes.Validator, 0, vals.Total)
		}

		for _, val := range respPb.Validators {
			vals.Validators = append(vals.Validators, convertValidator(val))
		}

		vals.Count += len(respPb.Validators)

		if len(respPb.Pagination.NextKey) == 0 {
			break
		}
	}

	return vals, nil
}

func convertValidator(c *tmservice.Validator) *cometbfttypes.Validator {
	pk := ed25519.PubKey(c.PubKey.Value)

	return &cometbfttypes.Validator{
		Address:          cometbfttypes.Address(c.Address),
		PubKey:           &pk,
		VotingPower:      c.VotingPower,
		ProposerPriority: c.ProposerPriority,
	}
}
