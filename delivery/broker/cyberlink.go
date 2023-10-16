package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishCyberLink(ctx context.Context, msg model.CyberLink) error {
	return b.marshalAndProduce(CyberLink, msg)
}

func (b *Broker) PublishCyberLinkMessage(ctx context.Context, msg model.CyberLinkMessage) error {
	return b.marshalAndProduce(CyberLinkMessage, msg)
}

func (b *Broker) PublishParticle(ctx context.Context, msg model.Particle) error {
	return b.marshalAndProduce(Particle, msg)
}
