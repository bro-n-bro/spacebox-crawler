package utils

import (
	grpcClient "bro-n-bro-osmosis/client/grpc"
	"bro-n-bro-osmosis/types"
	"context"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// RefreshBalance returns a function that when called refreshes the balance of the user having the given address
func RefreshBalance(address string, client banktypes.QueryClient) func() {
	return func() {
		var height int64
		//height, err := db.GetLastBlockHeight()
		//if err != nil {
		//	log.Error().Err(err).Str("module", "bank").
		//		Str("operation", "refresh balance").Msg("error while getting latest block height")
		//	return
		//}

		err := UpdateBalances([]string{address}, height, client)
		if err != nil {
			//log.Error().Err(err).Str("module", "bank").
			//	Str("operation", "refresh balance").Msg("error while updating balance")
		}
	}
}

// UpdateBalances updates the balances of the accounts having the given addresses,
// taking the data at the provided height
func UpdateBalances(addresses []string, height int64, bankClient banktypes.QueryClient) error {
	//log.Debug().Str("module", "bank").Int64("height", height).Msg("updating balances")
	header := grpcClient.GetHeightRequestHeader(height)

	var balances []types.AccountBalance
	for _, address := range addresses {
		balRes, err := bankClient.AllBalances(
			context.Background(),
			&banktypes.QueryAllBalancesRequest{Address: address},
			header,
		)
		if err != nil {
			return err
		}

		balances = append(balances, types.NewAccountBalance(
			address,
			balRes.Balances,
			height,
		))
	}

	// TODO:
	//err :=  db.SaveAccountBalances(balances)
	return nil
}
