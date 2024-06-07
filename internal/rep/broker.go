package rep

import (
	"context"
)

type Broker interface {
	// raw
	PublishRawBlock(ctx context.Context, b interface{}) error
	PublishRawTransaction(ctx context.Context, tx interface{}) error
	PublishRawBlockResults(ctx context.Context, br interface{}) error
	PublishRawGenesis(ctx context.Context, g interface{}) error
}
