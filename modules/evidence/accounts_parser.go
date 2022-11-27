package evidence

import (
	"bro-n-bro-osmosis/modules/messages"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
)

// EvidenceMessagesParser returns the list of all the accounts involved in the given
// message if it's related to the x/evidence module
func EvidenceMessagesParser(_ codec.Codec, cdkMsg sdk.Msg) ([]string, error) {
	switch msg := cdkMsg.(type) {

	case *evidencetypes.MsgSubmitEvidence:
		return []string{msg.Submitter}, nil

	}

	return nil, messages.MessageNotSupported(cdkMsg)
}
