package ibc

import (
	"context"
	"strings"
	"sync"

	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

const (
	denomPrefix = "ibc/"
)

// TODO: make better
var (
	lastHeight   int64
	lastHeightMu sync.RWMutex
)

func (m *Module) HandleBlock(ctx context.Context, block *types.Block) error {
	lastHeightMu.RLock()
	if lastHeight > block.Height {
		lastHeightMu.RUnlock()
		return nil
	}
	lastHeightMu.RUnlock()

	lastHeightMu.Lock()
	lastHeight = block.Height
	lastHeightMu.Unlock()

	var (
		nextKey     []byte
		denomHashes = make([]string, 0)
	)
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

		for _, coin := range respPb.Supply {
			if strings.HasPrefix(coin.Denom, denomPrefix) {
				hash := strings.TrimPrefix(coin.Denom, denomPrefix)
				m.denomCache.mu.RLock()
				if _, ok := m.denomCache.denomHashes[hash]; ok {
					m.denomCache.mu.RUnlock()
					continue
				}
				m.denomCache.mu.RUnlock()

				m.denomCache.mu.Lock()
				m.denomCache.denomHashes[hash] = struct{}{}
				m.denomCache.mu.Unlock()

				denomHashes = append(denomHashes, hash)
			}
		}

		nextKey = respPb.Pagination.NextKey
		if len(nextKey) == 0 {
			break
		}
	}

	return m.getAndPublishDenomTraces(ctx, denomHashes)
}

func (m *Module) getAndPublishDenomTraces(ctx context.Context, denomHashes []string) error {
	for _, hash := range denomHashes {
		resp, err := m.client.IbcTransferQueryClient.DenomTrace(
			ctx, &ibctransfertypes.QueryDenomTraceRequest{
				Hash: hash,
			},
		)
		if err != nil {
			return err
		}

		if resp.DenomTrace == nil {
			continue
		}

		if err = m.broker.PublishDenomTrace(ctx, model.DenomTrace{
			DenomHash: hash,
			Path:      resp.DenomTrace.Path,
			BaseDenom: resp.DenomTrace.BaseDenom,
		}); err != nil {
			return err
		}
	}

	return nil
}
