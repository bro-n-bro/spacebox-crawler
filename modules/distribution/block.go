package distribution

import (
	"context"
	"reflect"
	"sync"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	grpcClient "bro-n-bro-osmosis/client/grpc"
	"bro-n-bro-osmosis/modules/distribution/utils"
	"bro-n-bro-osmosis/types"
)

var (
	mu         sync.Mutex
	lastParams *types.DistributionParams
)

func (m *Module) HandleBlock(ctx context.Context, block *types.Block, _ *tmctypes.ResultValidators) error {
	// TODO: maybe use consensus client for get correct validators?
	go m.updateParams(block.Height)

	// Update the validator commissions
	go utils.UpdateValidatorsCommissionAmounts(block.Height, m.client.DistributionQueryClient)

	// Update the delegators commissions amounts
	go utils.UpdateDelegatorsRewardsAmounts(block.Height, m.client.DistributionQueryClient)
	return nil
}

func (m *Module) updateParams(height int64) {
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

	// TODO:
	params := types.NewDistributionParams(res.Params, height)
	mu.Lock()
	if lastParams == nil {
		m.log.Warn().Msg("set first params")
		lastParams = &params
	} else if !reflect.DeepEqual(lastParams.Params, params.Params) {
		m.log.Warn().
			Int64("last_height", lastParams.Height).
			Int64("cur_height", height).
			Str("last_params", lastParams.String()).
			Str("cur_params", params.String()).
			Msg("params not equal")
		lastParams = &params
	}
	mu.Unlock()

	//err = db.SaveDistributionParams(types.NewDistributionParams(res.Params, height))
	//if err != nil {
	//	log.Error().Str("module", "distribution").Err(err).
	//		Int64("height", height).
	//		Msg("error while saving params")
	//	return
	//}
}
