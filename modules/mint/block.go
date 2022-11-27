package mint

import (
	"context"

	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	grpcClient "bro-n-bro-osmosis/client/grpc"
	"bro-n-bro-osmosis/types"
)

func (m *Module) HandleBlock(ctx context.Context, block *types.Block, _ *tmctypes.ResultValidators) error {
	paramsResp, err := m.client.MintQueryClient.Params(
		ctx,
		&minttypes.QueryParamsRequest{},
		grpcClient.GetHeightRequestHeader(block.Height),
	)
	if err != nil {
		m.log.Error().Err(err).Int64("height", block.Height).Msg("error while getting params")
		return err
	}

	// not used in bdjuno
	// todo: call panic: invalid Go type types.Dec for field cosmos.mint.v1beta1.QueryInflationResponse.inflation
	//inflationResp, err := m.client.MintQueryClient.Inflation(
	//	ctx,
	//	&minttypes.QueryInflationRequest{},
	//	grpcClient.GetHeightRequestHeader(block.Height),
	//)
	//if err != nil {
	//	m.log.Error().Err(err).Int64("height", block.Height).Msg("error while getting inflation")
	//	return err
	//}
	//

	// TODO:

	_ = types.NewMintParams(paramsResp.Params, block.Height)

	return nil
}
