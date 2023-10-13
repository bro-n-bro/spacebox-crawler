package graph

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

type broker interface {
	PublishCyberLink(context.Context, model.CyberLink) error
	PublishCyberLinkMessage(context.Context, model.CyberLinkMessage) error
	PublishParticles(context.Context, model.Particles) error
}
