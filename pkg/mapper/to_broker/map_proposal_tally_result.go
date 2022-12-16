package to_broker

import (
	"github.com/hexy-dev/spacebox/broker/model"

	"bro-n-bro-osmosis/types"
)

func (tb ToBroker) MapProposalTallyResult(ptr types.TallyResult) model.ProposalTallyResult {
	return model.ProposalTallyResult{
		ProposalID: ptr.ProposalID,
		Yes:        ptr.Yes,
		Abstain:    ptr.Abstain,
		No:         ptr.No,
		NoWithVeto: ptr.NoWithVeto,
		Height:     ptr.Height,
	}
}
