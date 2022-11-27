package bank

import (
	"context"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	grpcClient "bro-n-bro-osmosis/client/grpc"
	"bro-n-bro-osmosis/types"
)

func (m *Module) HandleBlock(ctx context.Context, block *types.Block, _ *tmctypes.ResultValidators) error {
	params, err := m.getGovParams(ctx, block.Height)
	// TODO: UpdateProposal
	_ = params
	return err
}

func (m *Module) getGovParams(ctx context.Context, height int64) (*types.GovParams, error) {
	respDeposit, err := m.client.GovQueryClient.Params(
		ctx,
		&govtypes.QueryParamsRequest{ParamsType: govtypes.ParamDeposit},
		grpcClient.GetHeightRequestHeader(height),
	)
	if err != nil {
		return nil, err
	}

	respVoting, err := m.client.GovQueryClient.Params(
		ctx,
		&govtypes.QueryParamsRequest{ParamsType: govtypes.ParamVoting},
		grpcClient.GetHeightRequestHeader(height),
	)
	if err != nil {
		return nil, err

	}

	respTally, err := m.client.GovQueryClient.Params(
		ctx,
		&govtypes.QueryParamsRequest{ParamsType: govtypes.ParamTallying},
		grpcClient.GetHeightRequestHeader(height),
	)
	if err != nil {
		return nil, err
	}

	govParams := types.NewGovParams(
		types.NewVotingParams(respVoting.GetVotingParams()),
		types.NewDepositParam(respDeposit.GetDepositParams()),
		types.NewTallyParams(respTally.GetTallyParams()),
		height,
	)

	return govParams, nil
}
