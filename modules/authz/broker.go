package authz

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

type broker interface {
	PublishGrantMessage(context.Context, model.GrantMessage) error
	PublishRevokeMessage(context.Context, model.RevokeMessage) error
	PublishExecMessage(context.Context, model.ExecMessage) error
}
