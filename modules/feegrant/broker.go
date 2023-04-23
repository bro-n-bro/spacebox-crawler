package feegrant

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

type broker interface {
	PublishFeeAllowance(context.Context, model.FeeAllowance) error
	PublishGrantAllowanceMessage(context.Context, model.GrantAllowanceMessage) error
	PublishRevokeAllowanceMessage(context.Context, model.RevokeAllowanceMessage) error
}
