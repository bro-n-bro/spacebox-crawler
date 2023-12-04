package rep

import (
	"context"

	"github.com/bro-n-bro/spacebox-crawler/adapter/storage/model"
)

// Storage implementation needed for store some tmp data for correct processing
type Storage interface {
	GetBlockByHeight(ctx context.Context, height int64) (*model.Block, error)
	CreateBlock(ctx context.Context, block *model.Block) error
	SetProcessedStatus(ctx context.Context, height int64) error
	SetErrorStatus(ctx context.Context, height int64, msg string) error
	UpdateStatus(ctx context.Context, height int64, status model.Status) error
	GetErrorBlockHeights(ctx context.Context) ([]int64, error)

	InsertErrorTx(ctx context.Context, message model.Tx) error
	InsertErrorMessage(ctx context.Context, message model.Message) error

	Ping(ctx context.Context) error
}
