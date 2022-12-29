package slasing

import (
	"github.com/hexy-dev/spacebox-crawler/modules/messages"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
)

// SlashingMessagesParser returns the list of all the accounts involved in the given
// message if it's related to the x/slashing module
func SlashingMessagesParser(_ codec.Codec, sdkMsg sdk.Msg) ([]string, error) {
	// nolint:gocritic
	switch msg := sdkMsg.(type) {
	case *slashingtypes.MsgUnjail:
		return []string{msg.ValidatorAddr}, nil

	}

	return nil, messages.MessageNotSupported(sdkMsg)
}
