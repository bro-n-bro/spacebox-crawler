package distribution

import (
	"github.com/hexy-dev/spacebox-crawler/modules/messages"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// DistributionMessagesParser returns the list of all the accounts involved in the given
// message if it's related to the x/distribution module
func DistributionMessagesParser(_ codec.Codec, sdkMsg sdk.Msg) ([]string, error) {
	switch msg := sdkMsg.(type) {
	case *distrtypes.MsgSetWithdrawAddress:
		return []string{msg.DelegatorAddress, msg.WithdrawAddress}, nil

	case *distrtypes.MsgWithdrawDelegatorReward:
		return []string{msg.DelegatorAddress, msg.ValidatorAddress}, nil

	case *distrtypes.MsgWithdrawValidatorCommission:
		return []string{msg.ValidatorAddress}, nil

	case *distrtypes.MsgFundCommunityPool:
		return []string{msg.Depositor}, nil

	}

	return nil, messages.MessageNotSupported(sdkMsg)
}
