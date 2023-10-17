package slashing

import (
	"context"

	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

func (m *Module) HandleBlock(ctx context.Context, block *types.Block) error {
	paramsResp, err := m.client.SlashingQueryClient.Params(
		ctx,
		&slashingtypes.QueryParamsRequest{},
		grpcClient.GetHeightRequestHeader(block.Height),
	)
	if err != nil {
		m.log.Error().Err(err).Int64("height", block.Height).Msg("error while getting params")
		return err
	}

	if err = m.broker.PublishSlashingParams(ctx, model.SlashingParams{
		Height: block.Height,
		Params: model.RawSlashingParams{
			DowntimeJailDuration:    paramsResp.Params.DowntimeJailDuration,
			SignedBlocksWindow:      paramsResp.Params.SignedBlocksWindow,
			MinSignedPerWindow:      paramsResp.Params.MinSignedPerWindow.MustFloat64(),
			SlashFractionDoubleSign: paramsResp.Params.SlashFractionDoubleSign.MustFloat64(),
			SlashFractionDowntime:   paramsResp.Params.SlashFractionDowntime.MustFloat64(),
		},
	}); err != nil {
		return err
	}

	return nil
}
