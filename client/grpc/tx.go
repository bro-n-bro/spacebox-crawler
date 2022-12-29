package grpc

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"sync"

	"github.com/cosmos/cosmos-sdk/types/tx"
	types2 "github.com/tendermint/tendermint/types"
	"golang.org/x/sync/errgroup"
)

// Txs queries for all the transactions in a block. Transactions are returned
// to the sdk.TxResponse format which internally contains a sdk.Tx. An error is
// returned if any query fails.
func (c *Client) Txs(ctx context.Context, txs types2.Txs) ([]*tx.GetTxResponse, error) {
	mu := sync.Mutex{}
	txResponses := make([]*tx.GetTxResponse, 0, len(txs))

	g, ctx := errgroup.WithContext(ctx)
	for _, tmTx := range txs {
		tmTx := tmTx
		g.Go(func() error {
			respPb, err := c.TxService.GetTx(ctx, &tx.GetTxRequest{Hash: fmt.Sprintf("%X", tmTx.Hash())})
			if err != nil {
				return err
			}

			mu.Lock()
			txResponses = append(txResponses, &tx.GetTxResponse{Tx: respPb.Tx, TxResponse: respPb.TxResponse})
			mu.Unlock()
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return txResponses, nil
}

func (c *Client) TxsOld(ctx context.Context, txs types2.Txs) ([]*tx.GetTxResponse, error) {
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
