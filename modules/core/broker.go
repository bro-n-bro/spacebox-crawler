package core

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

type broker interface {
	PublishBlock(context.Context, model.Block) error
	PublishMessage(ctx context.Context, message model.Message) error
	PublishTransaction(ctx context.Context, tx model.Transaction) error
	PublishValidatorPreCommit(ctx context.Context, vp model.ValidatorPreCommit) error
	PublishValidatorVotingPower(_ context.Context, vp model.ValidatorVotingPower) error
}
