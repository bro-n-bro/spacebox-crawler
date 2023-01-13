package grpc

import (
	"context"
	"encoding/hex"
	"log"

	"github.com/cosmos/cosmos-sdk/types/tx"
	types2 "github.com/tendermint/tendermint/types"
)

// Txs queries for all the transactions in a block. Transactions are returned
// in sdk.TxResponse format which internally contains a sdk.Tx. An error is
// returned if any query fails.
func (c *Client) Txs(ctx context.Context, txs types2.Txs) ([]*tx.GetTxResponse, error) {
	txResponses := make([]*tx.GetTxResponse, 0, len(txs))

	for _, tmTx := range txs {
		respPb, err := c.TxService.GetTx(ctx, &tx.GetTxRequest{Hash: hex.EncodeToString(tmTx.Hash())})
		if err != nil {
			log.Println("GetTx error:", err)
			continue
		}

		txResponses = append(txResponses, &tx.GetTxResponse{Tx: respPb.Tx, TxResponse: respPb.TxResponse})
	}

	return txResponses, nil
}
