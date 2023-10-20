package grid

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

type broker interface {
	PublishGridParams(ctx context.Context, mp model.GridParams) error

	PublishCreateRouteMessage(ctx context.Context, msg model.CreateRouteMessage) error
	PublishEditRouteMessage(ctx context.Context, msg model.EditRouteMessage) error
	PublishEditRouteNameMessage(ctx context.Context, msg model.EditRouteNameMessage) error
	PublishDeleteRouteMessage(ctx context.Context, msg model.DeleteRouteMessage) error
}
