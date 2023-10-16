package dmn

import (
	"context"

	dmntypes "github.com/cybercongress/go-cyber/x/dmn/types"
	"github.com/pkg/errors"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

func (m *Module) HandleBlock(ctx context.Context, block *types.Block) error {
	paramsResp, err := m.client.DMNQueryClient.Params(
		ctx,
		&dmntypes.QueryParamsRequest{},
		grpcClient.GetHeightRequestHeader(block.Height),
	)
	if err != nil {
		m.log.Error().Err(err).Int64("height", block.Height).Msg("error while getting params")
		return err
	}

	if err = m.broker.PublishDMNParams(ctx, model.DMNParams{
		Height: block.Height,
		Params: model.DMParams{
			MaxSlots: int64(paramsResp.Params.MaxSlots),
			MaxGas:   int64(paramsResp.Params.MaxGas),
			FeeTTL:   int64(paramsResp.Params.FeeTtl),
		},
	}); err != nil {
		return errors.Wrap(err, "PublishDMNParams error")
	}

	return nil
}
