package rep

import (
	"context"

	cometbftcoretypes "github.com/cometbft/cometbft/rpc/core/types"
	cometbfttypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/types/tx"

	"github.com/bro-n-bro/spacebox-crawler/types"
)

type (
	GrpcClient interface {
		Block(ctx context.Context, height int64) (*cometbftcoretypes.ResultBlock, error)
		Validators(ctx context.Context, height int64) (*cometbftcoretypes.ResultValidators, error)

		Txs(ctx context.Context, txs cometbfttypes.Txs) ([]*tx.GetTxResponse, error)
	}

	RPCClient interface {
		WsEnabled() bool

		SubscribeNewBlocks(ctx context.Context) (<-chan cometbftcoretypes.ResultEvent, error)
		Genesis(ctx context.Context) (*cometbfttypes.GenesisDoc, error)
		GetLastBlockHeight(ctx context.Context) (int64, error)
		GetBlockEvents(ctx context.Context, height int64) (begin, end types.BlockerEvents, err error)
	}
)
