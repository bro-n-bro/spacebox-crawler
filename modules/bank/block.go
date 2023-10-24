package bank

import (
	"context"

	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

func (m *Module) HandleBlock(ctx context.Context, block *types.Block) error {
	supply := model.Supply{
		Height: block.Height,
	}

	var nextKey []byte

	for {
		respPb, err := m.client.BankQueryClient.TotalSupply(
			ctx,
			&banktypes.QueryTotalSupplyRequest{
				Pagination: &query.PageRequest{
					Key:        nextKey,
					Limit:      100,
					CountTotal: true,
				},
			},
			grpcClient.GetHeightRequestHeader(block.Height))
		if err != nil {
			return err
		}

		if cap(supply.Coins) == 0 {
			supply.Coins = make(model.Coins, 0, respPb.Pagination.Total)
		}

		supply.Coins = append(supply.Coins, m.tbM.MapCoins(types.NewCoinsFromSDK(respPb.Supply))...)

		nextKey = respPb.Pagination.NextKey
		if len(nextKey) == 0 {
			break
		}
	}

	// TODO: test it
	return m.broker.PublishSupply(ctx, supply)
}
