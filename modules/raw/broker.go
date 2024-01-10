package raw

import "context"

type broker interface {
	PublishRawBlock(ctx context.Context, b interface{}) error
	PublishRawTransaction(ctx context.Context, tx interface{}) error
	PublishRawBlockResults(ctx context.Context, br interface{}) error
}
