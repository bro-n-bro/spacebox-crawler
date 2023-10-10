package graph

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

type broker interface {
	PublishCyberlinkMessage(context.Context, model.CyberlinkMessage) error
}
