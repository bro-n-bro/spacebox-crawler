package bank

import (
	"bro-n-bro-osmosis/modules/messages"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// BankAccountsParser returns the list of all the accounts involved in the given
// message if it's related to the x/bank module
func BankAccountsParser(_ codec.Codec, sdkMsg sdk.Msg) ([]string, error) {
	switch msg := sdkMsg.(type) {
	case *banktypes.MsgSend:
		return []string{msg.ToAddress, msg.FromAddress}, nil

	case *banktypes.MsgMultiSend:
		var addresses []string
		for _, i := range msg.Inputs {
			addresses = append(addresses, i.Address)
		}
		for _, o := range msg.Outputs {
			addresses = append(addresses, o.Address)
		}
		return addresses, nil
	}

	return nil, messages.MessageNotSupported(sdkMsg)
}
