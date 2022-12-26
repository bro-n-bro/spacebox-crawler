package rep

import (
	"context"

	"bro-n-bro-osmosis/adapter/storage/model"
)

// Storage implementation needed for store some tmp data for correct processing
type Storage interface {
	HasBlock(ctx context.Context, height int64) (bool, error)
	GetBlockStatus(ctx context.Context, height int64) (model.Status, error)
	CreateBlock(ctx context.Context, block *model.Block) error
	SetProcessedStatus(ctx context.Context, height int64) error
	SetErrorStatus(ctx context.Context, height int64) error
	UpdateStatus(ctx context.Context, height int64, status model.Status) error
	GetErrorBlockHeights(ctx context.Context) ([]int64, error)
}
