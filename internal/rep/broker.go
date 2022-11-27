package rep

import (
	"context"

	"bro-n-bro-osmosis/adapter/broker/model"
)

type Broker interface {
	PublishAccounts(context.Context, []model.Account) error
	PublishBlock(context.Context, model.Block) error
	PublishSupply(context.Context, model.Supply) error
	PublishSendMessage(context.Context, model.SendMessage) error
	PublishDelegationRewardMessage(context.Context, model.DelegationRewardMessage) error
	PublishProposalVoteMessage(context.Context, model.ProposalVoteMessage) error
	PublishRedelegationMessage(context.Context, model.RedelegationMessage) error
	PublishUnbondingDelegationMessage(context.Context, model.UnbondingDelegationMessage) error
}
