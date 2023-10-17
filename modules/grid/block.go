package grid

import (
	"context"

	gridtypes "github.com/cybercongress/go-cyber/x/grid/types"
	"github.com/pkg/errors"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

func (m *Module) HandleBlock(ctx context.Context, block *types.Block) error {
	paramsResp, err := m.client.GridQueryClient.Params(
		ctx,
		&gridtypes.QueryParamsRequest{},
		grpcClient.GetHeightRequestHeader(block.Height),
	)
	if err != nil {
		m.log.Error().Err(err).Int64("height", block.Height).Msg("error while getting params")
		return err
	}

	if err = m.broker.PublishGridParams(ctx, model.GridParams{
		Height: block.Height,
		Params: model.RawGridParams{
			MaxRoutes: int64(paramsResp.Params.MaxRoutes),
		},
	}); err != nil {
		return errors.Wrap(err, "PublishGridParams error")
	}

	return nil
}
