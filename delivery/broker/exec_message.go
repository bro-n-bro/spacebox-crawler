package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishExecMessage(ctx context.Context, em model.ExecMessage) error {
	return b.marshalAndProduce(ExecMessage, em)
}
