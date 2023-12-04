package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishCyberlink(ctx context.Context, msg model.Cyberlink) error {
	return b.marshalAndProduce(Cyberlink, msg)
}

func (b *Broker) PublishCyberlinkMessage(ctx context.Context, msg model.CyberlinkMessage) error {
	return b.marshalAndProduce(CyberlinkMessage, msg)
}

func (b *Broker) PublishParticle(ctx context.Context, msg model.Particle) error {
	return b.marshalAndProduce(Particle, msg)
}
