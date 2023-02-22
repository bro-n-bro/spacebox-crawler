package auth

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

type broker interface {
	PublishAccount(context.Context, model.Account) error
}
