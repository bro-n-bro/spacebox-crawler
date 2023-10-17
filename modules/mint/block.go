package mint

import (
	"context"

	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/pkg/errors"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

func (m *Module) HandleBlock(ctx context.Context, block *types.Block) error {
	paramsResp, err := m.client.MintQueryClient.Params(
		ctx,
		&minttypes.QueryParamsRequest{},
		grpcClient.GetHeightRequestHeader(block.Height),
	)
	if err != nil {
		m.log.Error().Err(err).Int64("height", block.Height).Msg("error while getting params")
		return err
	}

	// todo: call panic: invalid Go type types.Dec for field cosmos.mint.v1beta1.QueryInflationResponse.inflation
	// inflationResp, err := m.client.MintQueryClient.Inflation(
	//	ctx,
	//	&minttypes.QueryInflationRequest{},
	//	grpcClient.GetHeightRequestHeader(block.Height),
	// )
	// if err != nil {
	//	m.log.Error().Err(err).Int64("height", block.Height).Msg("error while getting inflation")
	//	return err
	// }
	// _ = inflationResp

	// m.client.MintQueryClient.AnnualProvisions()

	// TODO: maybe check diff from mongo in my side?
	// TODO: test it
	if err = m.broker.PublishMintParams(ctx, model.MintParams{
		Height: block.Height,
		Params: model.RawMintParams{
			MintDenom:           paramsResp.Params.MintDenom,
			InflationRateChange: paramsResp.Params.InflationRateChange.MustFloat64(),
			InflationMax:        paramsResp.Params.InflationMax.MustFloat64(),
			InflationMin:        paramsResp.Params.InflationMin.MustFloat64(),
			GoalBonded:          paramsResp.Params.GoalBonded.MustFloat64(),
			BlocksPerYear:       paramsResp.Params.BlocksPerYear,
		},
	}); err != nil {
		return errors.Wrap(err, "PublishMintParams error")
	}

	// TODO: got a panic: invalid Go type types.Dec for field
	// cosmos.mint.v1beta1.QueryAnnualProvisionsResponse.annual_provisions

	// annualProvResp, err := m.client.MintQueryClient.AnnualProvisions(
	//	ctx,
	//	&minttypes.QueryAnnualProvisionsRequest{},
	//	grpcClient.GetHeightRequestHeader(block.Height),
	// )
	// if err != nil {
	//	m.log.Error().Err(err).Int64("height", block.Height).Msg("error while annual provision")
	//	return err
	// }

	// var annualProvision float64
	// if !annualProvResp.AnnualProvisions.IsNil() {
	//	annualProvision = annualProvResp.AnnualProvisions.MustFloat64()
	// }

	// if err := m.broker.PublishAnnualProvision(ctx,
	// m.tbM.MapAnnualProvision(block.Height, annualProvision)); err != nil {
	//	return errors.Wrap(err, "PublishAnnualProvision error")
	// }

	return nil
}
