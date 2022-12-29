package slasing

import (
	"context"

	"github.com/cosmos/cosmos-sdk/types/query"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	grpcClient "github.com/hexy-dev/spacebox-crawler/client/grpc"
	"github.com/hexy-dev/spacebox-crawler/types"
)

func (m *Module) HandleBlock(ctx context.Context, block *types.Block, _ *tmctypes.ResultValidators) error {
	// Update the signing infos
	err := m.updateSigningInfo(block.Height)
	if err != nil {
		m.log.Error().Int64("height", block.Height).
			Err(err).Msg("error while updating signing info")
	}

	err = m.updateSlashingParams(block.Height)
	if err != nil {
		m.log.Error().Int64("height", block.Height).Err(err).Msg("error while updating params")
	}

	return nil
}

// updateSigningInfo reads from the LCD the current staking pool and stores its value inside the database
func (m *Module) updateSigningInfo(height int64) error {

	signingInfos, err := m.getSigningInfos(height)
	if err != nil {
		return err
	}

	// TODO:
	_ = signingInfos
	return nil
}

// updateSlashingParams gets the slashing params for the given height, and stores them inside the database
func (m *Module) updateSlashingParams(height int64) error {

	res, err := m.client.SlashingQueryClient.Params(
		context.Background(),
		&slashingtypes.QueryParamsRequest{},
		grpcClient.GetHeightRequestHeader(height),
	)
	if err != nil {
		return err
	}

	// TODO:
	_ = types.NewSlashingParams(res.Params, height)
	return nil
}

func (m *Module) getSigningInfos(height int64) ([]types.ValidatorSigningInfo, error) {
	var signingInfos []slashingtypes.ValidatorSigningInfo

	header := grpcClient.GetHeightRequestHeader(height)

	var nextKey []byte
	var stop = false
	for !stop {
		res, err := m.client.SlashingQueryClient.SigningInfos(
			context.Background(),
			&slashingtypes.QuerySigningInfosRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 1000, // Query 1000 signing infos at a time
				},
			},
			header,
		)
		if err != nil {
			return nil, err
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		signingInfos = append(signingInfos, res.Info...)
	}

	infos := make([]types.ValidatorSigningInfo, len(signingInfos))
	for index, info := range signingInfos {
		infos[index] = types.NewValidatorSigningInfo(
			info.Address,
			height,
			info.StartHeight,
			info.IndexOffset,
			info.MissedBlocksCounter,
			info.JailedUntil,
			info.Tombstoned,
		)
	}
	return infos, nil
}
