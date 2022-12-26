package mint

import (
	"context"

	"github.com/pkg/errors"

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
	//_ = inflationResp

	// m.client.MintQueryClient.AnnualProvisions()

	// TODO: maybe check diff from mongo in my side?
	params := types.NewMintParams(paramsResp.Params, block.Height)
	// TODO: test it
	err = m.broker.PublishMintParams(ctx, m.tbM.MapMingParams(params))
	if err != nil {
		return errors.Wrap(err, "PublishMintParams error")
	}

	// TODO: go panic: invalid Go type types.Dec for field cosmos.mint.v1beta1.QueryAnnualProvisionsResponse.annual_provisions
	//annualProvResp, err := m.client.MintQueryClient.AnnualProvisions(
	//	ctx,
	//	&minttypes.QueryAnnualProvisionsRequest{},
	//	grpcClient.GetHeightRequestHeader(block.Height),
	//)
	//if err != nil {
	//	m.log.Error().Err(err).Int64("height", block.Height).Msg("error while annual provision")
	//	return err
	//}

	//var annualProvision float64
	//if !annualProvResp.AnnualProvisions.IsNil() {
	//	annualProvision = annualProvResp.AnnualProvisions.MustFloat64()
	//}

	//if err := m.broker.PublishAnnualProvision(ctx, m.tbM.MapAnnualProvision(block.Height, annualProvision)); err != nil {
	//	return errors.Wrap(err, "PublishAnnualProvision error")
	//}

	return nil
}
