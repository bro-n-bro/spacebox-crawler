package slashing

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

type broker interface {
	PublishSlashingParams(ctx context.Context, params model.SlashingParams) error
	PublishUnJailMessage(ctx context.Context, msg model.UnjailMessage) error
	PublishHandleValidatorSignature(ctx context.Context, msg model.HandleValidatorSignature) error
}
