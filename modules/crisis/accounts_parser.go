package crisis

import (
	"bro-n-bro-osmosis/modules/messages"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
)

// CrisisMessagesParser returns the list of all the accounts involved in the given
// message if it's related to the x/crisis module
func CrisisMessagesParser(_ codec.Codec, sdkMsg sdk.Msg) ([]string, error) {
	switch msg := sdkMsg.(type) {

	case *crisistypes.MsgVerifyInvariant:
		return []string{msg.Sender}, nil
	}

	return nil, messages.MessageNotSupported(sdkMsg)
}
