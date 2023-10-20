package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishGridParams(ctx context.Context, msg model.GridParams) error {
	return b.marshalAndProduce(GridParams, msg)
}

func (b *Broker) PublishCreateRouteMessage(ctx context.Context, msg model.CreateRouteMessage) error {
	return b.marshalAndProduce(CreateRouteMessage, msg)
}

func (b *Broker) PublishEditRouteMessage(ctx context.Context, msg model.EditRouteMessage) error {
	return b.marshalAndProduce(EditRouteMessage, msg)
}

func (b *Broker) PublishEditRouteNameMessage(ctx context.Context, msg model.EditRouteNameMessage) error {
	return b.marshalAndProduce(EditRouteNameMessage, msg)
}

func (b *Broker) PublishDeleteRouteMessage(ctx context.Context, msg model.DeleteRouteMessage) error {
	return b.marshalAndProduce(DeleteRouteMessage, msg)
}
