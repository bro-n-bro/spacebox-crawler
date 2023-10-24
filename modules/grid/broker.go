package grid

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

type broker interface {
	PublishGridParams(context.Context, model.GridParams) error

	PublishRoute(context.Context, model.Route) error
	PublishCreateRouteMessage(context.Context, model.CreateRouteMessage) error
	PublishEditRouteMessage(context.Context, model.EditRouteMessage) error
	PublishEditRouteNameMessage(context.Context, model.EditRouteNameMessage) error
	PublishDeleteRouteMessage(context.Context, model.DeleteRouteMessage) error
}
