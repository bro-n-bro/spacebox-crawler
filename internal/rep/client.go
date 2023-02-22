package rep

import (
	"context"

	"github.com/cosmos/cosmos-sdk/types/tx"
	tmcoretypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

type (
	GrpcClient interface {
		Block(ctx context.Context, height int64) (*tmcoretypes.ResultBlock, error)
		Validators(ctx context.Context, height int64) (*tmcoretypes.ResultValidators, error)

		Txs(ctx context.Context, txs tmtypes.Txs) ([]*tx.GetTxResponse, error)
	}

	RPCClient interface {
		WsEnabled() bool

		SubscribeNewBlocks(ctx context.Context) (<-chan tmcoretypes.ResultEvent, error)
		Genesis(ctx context.Context) (*tmtypes.GenesisDoc, error)
		GetLastBlockHeight(ctx context.Context) (int64, error)
	}
)
