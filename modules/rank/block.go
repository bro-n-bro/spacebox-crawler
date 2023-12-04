package rank

import (
	"context"

	ranktypes "github.com/cybercongress/go-cyber/x/rank/types"
	"github.com/pkg/errors"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

func (m *Module) HandleBlock(ctx context.Context, block *types.Block) error {
	paramsResp, err := m.client.RankQueryClient.Params(
		ctx,
		&ranktypes.QueryParamsRequest{},
		grpcClient.GetHeightRequestHeader(block.Height),
	)
	if err != nil {
		m.log.Error().Err(err).Int64("height", block.Height).Msg("error while getting params")
		return err
	}

	if err = m.broker.PublishRankParams(ctx, model.RankParams{
		Height: block.Height,
		Params: model.RawRankParams{
			CalculationPeriod: paramsResp.Params.CalculationPeriod,
			DampingFactor:     paramsResp.Params.DampingFactor.MustFloat64(),
			Tolerance:         paramsResp.Params.Tolerance.MustFloat64(),
		},
	}); err != nil {
		return errors.Wrap(err, "PublishRankParams error")
	}

	return nil
}
