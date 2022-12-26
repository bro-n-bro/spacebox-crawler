package bank

import (
	"bro-n-bro-osmosis/modules/messages"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// GovMessagesParser returns the list of all the accounts involved in the given
// message if it's related to the x/gov module
func GovMessagesParser(cdc codec.Codec, sdkMsg sdk.Msg) ([]string, error) {

	switch msg := sdkMsg.(type) {
	case *govtypes.MsgSubmitProposal:
		addresses := []string{msg.Proposer}

		var content govtypes.Content
		err := cdc.UnpackAny(msg.Content, &content)
		if err != nil {
			return nil, err
		}

		// Get addresses from contents
		switch content := content.(type) {
		case *distrtypes.CommunityPoolSpendProposal:
			addresses = append(addresses, content.Recipient)
		}

		return addresses, nil

	case *govtypes.MsgDeposit:
		return []string{msg.Depositor}, nil

	case *govtypes.MsgVote:
		return []string{msg.Voter}, nil

	}

	return nil, messages.MessageNotSupported(sdkMsg)
}