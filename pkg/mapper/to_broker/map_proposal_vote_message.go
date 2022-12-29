package to_broker

import (
	"github.com/hexy-dev/spacebox/broker/model"

	"github.com/hexy-dev/spacebox-crawler/types"
)

func (tb ToBroker) MapProposalVoteMessage(pvm types.ProposalVoteMessage) model.ProposalVoteMessage {
	return model.ProposalVoteMessage{
		ProposalID: pvm.ProposalID,
		Voter:      pvm.Voter,
		Option:     pvm.Option.String(),
		Height:     pvm.Height,
	}
}
