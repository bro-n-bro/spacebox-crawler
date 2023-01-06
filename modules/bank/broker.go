package bank

import (
	"context"

	"github.com/hexy-dev/spacebox/broker/model"
)

type broker interface {
	PublishSupply(context.Context, model.Supply) error
	PublishMultiSendMessage(ctx context.Context, msm model.MultiSendMessage) error
	PublishSendMessage(context.Context, model.SendMessage) error
	PublishAccountBalance(ctx context.Context, ab model.AccountBalance) error
}
