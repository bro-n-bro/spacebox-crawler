package wasm

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

type broker interface {
	PublishCyberlink(context.Context, model.Cyberlink) error
	PublishParticle(context.Context, model.Particle) error
}
