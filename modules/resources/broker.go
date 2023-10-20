package resources

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

type broker interface {
	PublishInvestmintMessage(ctx context.Context, mp model.InvestmintMessage) error
}
