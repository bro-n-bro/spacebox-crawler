package auth

import (
	"context"

	"github.com/hexy-dev/spacebox/broker/model"
)

type broker interface {
	PublishAccount(context.Context, model.Account) error
}
