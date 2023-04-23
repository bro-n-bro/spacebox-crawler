package authz

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

// Broker is an interface for publishing authz messages.
type broker interface {
	PublishAuthzGrant(context.Context, model.AuthzGrant) error
	PublishGrantMessage(context.Context, model.GrantMessage) error
	PublishRevokeMessage(context.Context, model.RevokeMessage) error
	PublishExecMessage(context.Context, model.ExecMessage) error
}
