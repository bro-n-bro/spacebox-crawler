package auth

import (
	"context"

	"github.com/hexy-dev/spacebox/broker/model"
)

type broker interface {
	PublishAccounts(context.Context, []model.Account) error
}
