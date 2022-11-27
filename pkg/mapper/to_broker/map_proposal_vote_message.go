package to_broker

import (
	"bro-n-bro-osmosis/adapter/broker/model"
	"bro-n-bro-osmosis/types"
)

func (tb ToBroker) MapProposalVoteMessage(pvm types.ProposalVoteMessage) model.ProposalVoteMessage {
	return model.ProposalVoteMessage{
		ProposalID: pvm.ProposalID,
		Voter:      pvm.Voter,
		Option:     pvm.Option.String(),
		Height:     pvm.Height,
	}
}
