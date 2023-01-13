package bank

import (
	"context"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	grpcClient "github.com/hexy-dev/spacebox-crawler/client/grpc"
	"github.com/hexy-dev/spacebox-crawler/types"
	"github.com/hexy-dev/spacebox/broker/model"
)

func (m *Module) updateBalance(ctx context.Context, addresses []string, height int64) error {
	header := grpcClient.GetHeightRequestHeader(height)

	for _, address := range addresses {
		balRes, err := m.client.BankQueryClient.AllBalances(
			context.Background(),
			&banktypes.QueryAllBalancesRequest{Address: address},
			header,
		)
		if err != nil {
			return err
		}

		if err = m.broker.PublishAccountBalance(
			ctx,
			model.AccountBalance{
				Address: address,
				Height:  height,
				Coins:   m.tbM.MapCoins(types.NewCoinsFromCdk(balRes.Balances)),
			}); err != nil {
			return err
		}
	}

	return nil
}
