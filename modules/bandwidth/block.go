package bandwidth

import (
	"context"

	bandwidthtypes "github.com/cybercongress/go-cyber/x/bandwidth/types"
	"github.com/pkg/errors"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

func (m *Module) HandleBlock(ctx context.Context, block *types.Block) error {
	paramsResp, err := m.client.BandwidthQueryClient.Params(
		ctx,
		&bandwidthtypes.QueryParamsRequest{},
		grpcClient.GetHeightRequestHeader(block.Height),
	)
	if err != nil {
		m.log.Error().Err(err).Int64("height", block.Height).Msg("error while getting params")
		return err
	}

	if err = m.broker.PublishBandwidthParams(ctx, model.BandwidthParams{
		Height: block.Height,
		Params: model.RawBandwidthParams{
			RecoveryPeriod:    paramsResp.Params.RecoveryPeriod,
			AdjustPricePeriod: paramsResp.Params.AdjustPricePeriod,
			BasePrice:         paramsResp.Params.BasePrice.MustFloat64(),
			BaseLoad:          paramsResp.Params.BaseLoad.MustFloat64(),
			MaxBlockBandwidth: paramsResp.Params.MaxBlockBandwidth,
		},
	}); err != nil {
		return errors.Wrap(err, "PublishBandwidthParams error")
	}

	return nil
}
