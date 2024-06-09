package rep

import (
	"context"

	"github.com/cosmos/cosmos-sdk/types/tx"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/bro-n-bro/spacebox-crawler/v2/types"
)

type (
	GrpcClient interface {
		Block(ctx context.Context, height int64) (*coretypes.ResultBlock, error)
		Validators(ctx context.Context, height int64) (*coretypes.ResultValidators, error)

		Txs(ctx context.Context, height int64, txs tmtypes.Txs) ([]*tx.GetTxResponse, error)
	}

	RPCClient interface {
		SubscribeNewBlocks(ctx context.Context) (<-chan coretypes.ResultEvent, error)
		Genesis(ctx context.Context) (*tmtypes.GenesisDoc, error)
		GetLastBlockHeight(ctx context.Context) (int64, error)
		GetBlockEvents(ctx context.Context, height int64) (begin, end types.BlockerEvents, err error)
	}
)
