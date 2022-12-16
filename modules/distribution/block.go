package distribution

import (
	"context"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	grpcClient "bro-n-bro-osmosis/client/grpc"
	"bro-n-bro-osmosis/modules/distribution/utils"
	"bro-n-bro-osmosis/types"
)

func (m *Module) HandleBlock(ctx context.Context, block *types.Block, vals *tmctypes.ResultValidators) error {
	// TODO: maybe use consensus client for get correct validators?
	go m.updateParams(ctx, block.Height)

	// Update the validator commissions
	go utils.UpdateValidatorsCommissionAmounts(block.Height, m.client.DistributionQueryClient)

	// Update the delegators commissions amounts
	go utils.UpdateDelegatorsRewardsAmounts(block.Height, m.client.DistributionQueryClient)

	// TODO: client.community pull
	return nil
}

func (m *Module) updateParams(ctx context.Context, height int64) {
	//log.Debug().Str("module", "distribution").Int64("height", height).
	//	Msg("updating params")

	res, err := m.client.DistributionQueryClient.Params(
		context.Background(),
		&distrtypes.QueryParamsRequest{},
		grpcClient.GetHeightRequestHeader(height),
	)
	if err != nil {
		//log.Error().Str("module", "distribution").Err(err).
		//	Int64("height", height).
		//	Msg("error while getting params")
		return
	}

	// TODO: maybe check diff from mongo in my side?
	params := types.NewDistributionParams(res.Params, height)
	// TODO: test it
	if err := m.broker.PublishDistributionParams(ctx, m.tbM.MapDistributionParams(params)); err != nil {
		m.log.Error().Int64("height", height).Err(err).Msg("PublishDistributionParams error")
	}
}
