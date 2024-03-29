package core

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	feegranttypes "github.com/cosmos/cosmos-sdk/x/feegrant"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	graph "github.com/cybercongress/go-cyber/x/graph/types"
)

var (
	// CosmosMessageAddressesParser represents a MsgAddrParser that parses a
	// Cosmos message and returns all the involved addresses (both accounts and validators)
	CosmosMessageAddressesParser = JoinMessageParsers(
		BankMessagesParser,
		CrisisMessagesParser,
		DistributionMessagesParser,
		EvidenceMessagesParser,
		GovMessagesParser,
		IBCTransferMessagesParser,
		SlashingMessagesParser,
		StakingMessagesParser,
		FeeGrantMessagesParser,
		AuthzMessagesParser,
		GraphMessagesParser,

		DefaultMessagesParser,
	)
)

type (
	// MsgAddrParser represents a function that extracts all the
	// involved addresses from a provided message (both accounts and validators)
	MsgAddrParser = func(cdc codec.Codec, msg sdk.Msg) []string
)

// JoinMessageParsers joins together all the given parsers, calling them in order
func JoinMessageParsers(parsers ...MsgAddrParser) MsgAddrParser {
	return func(cdc codec.Codec, msg sdk.Msg) []string {
		// https://github.com/bro-n-bro/spacebox-crawler/issues/131
		if msg == nil {
			return nil
		}

		for _, parser := range parsers {
			// Try getting the addresses

			// If some addresses are found, return them
			if addresses := parser(cdc, msg); len(addresses) > 0 {
				return addresses
			}
		}

		return nil
	}
}

// DefaultMessagesParser represents the default messages parser that simply returns the list
// of all the signers of a message
func DefaultMessagesParser(_ codec.Codec, incomingMsg sdk.Msg) []string {
	var (
		cosmosSigners = incomingMsg.GetSigners()
		signers       = make([]string, len(cosmosSigners))
	)

	for index, signer := range cosmosSigners {
		signers[index] = signer.String()
	}

	return signers
}

// BankMessagesParser returns the list of all the accounts involved in the given
// message if it's related to the x/bank module
func BankMessagesParser(_ codec.Codec, incomingMsg sdk.Msg) []string {
	switch msg := incomingMsg.(type) {
	case *banktypes.MsgSend:
		return []string{msg.ToAddress, msg.FromAddress}
	case *banktypes.MsgMultiSend:
		var addresses []string

		for _, i := range msg.Inputs {
			addresses = append(addresses, i.Address)
		}

		for _, o := range msg.Outputs {
			addresses = append(addresses, o.Address)
		}

		return addresses
	}

	return nil
}

// CrisisMessagesParser returns the list of all the accounts involved in the given
// message if it's related to the x/crisis module
func CrisisMessagesParser(_ codec.Codec, incomingMsg sdk.Msg) []string {
	// nolint:gocritic
	switch msg := incomingMsg.(type) {
	case *crisistypes.MsgVerifyInvariant:
		return []string{msg.Sender}
	}

	return nil
}

// DistributionMessagesParser returns the list of all the accounts involved in the given
// message if it's related to the x/distribution module
func DistributionMessagesParser(_ codec.Codec, incomingMsg sdk.Msg) []string {
	switch msg := incomingMsg.(type) {
	case *distrtypes.MsgSetWithdrawAddress:
		return []string{msg.DelegatorAddress, msg.WithdrawAddress}
	case *distrtypes.MsgWithdrawDelegatorReward:
		return []string{msg.DelegatorAddress, msg.ValidatorAddress}
	case *distrtypes.MsgWithdrawValidatorCommission:
		return []string{msg.ValidatorAddress}
	case *distrtypes.MsgFundCommunityPool:
		return []string{msg.Depositor}
	}

	return nil
}

// EvidenceMessagesParser returns the list of all the accounts involved in the given
// message if it's related to the x/evidence module
func EvidenceMessagesParser(_ codec.Codec, incomingMsg sdk.Msg) []string {
	// nolint:gocritic
	switch msg := incomingMsg.(type) {
	case *evidencetypes.MsgSubmitEvidence:
		return []string{msg.Submitter}
	}

	return nil
}

// GovMessagesParser returns the list of all the accounts involved in the given
// message if it's related to the x/gov module
func GovMessagesParser(cdc codec.Codec, incomingMsg sdk.Msg) []string {
	switch msg := incomingMsg.(type) {
	case *govtypes.MsgSubmitProposal:
		var (
			addresses = []string{msg.Proposer}
			content   govtypes.Content
		)

		if err := cdc.UnpackAny(msg.Content, &content); err != nil {
			return nil
		}

		//nolint:gocritic,staticcheck
		// Get addresses from contents
		switch content := content.(type) {
		case *distrtypes.CommunityPoolSpendProposal:
			addresses = append(addresses, content.Recipient)
		}

		return addresses
	case *govtypes.MsgDeposit:
		return []string{msg.Depositor}
	case *govtypes.MsgVote:
		return []string{msg.Voter}
	}

	return nil
}

//
// // IBCTransferMessagesParser returns the list of all the accounts involved in the given
// // message if it's related to the x/iBCTransfer module
// func IBCTransferMessagesParser(_ codec.Codec, cosmosMsg sdk.Msg) ([]string, error) {
//	switch msg := cosmosMsg.(type) {
//
//	case *ibctransfertypes.MsgTransfer:
//		return []string{msg.Sender, msg.Receiver}, nil
//
//	}
//
//	return nil, MessageNotSupported(cosmosMsg)
// }

// SlashingMessagesParser returns the list of all the accounts involved in the given
// message if it's related to the x/slashing module
func SlashingMessagesParser(_ codec.Codec, incomingMsg sdk.Msg) []string {
	// nolint:gocritic
	switch msg := incomingMsg.(type) {
	case *slashingtypes.MsgUnjail:
		return []string{msg.ValidatorAddr}
	}

	return nil
}

// StakingMessagesParser returns the list of all the accounts involved in the given
// message if it's related to the x/staking module
func StakingMessagesParser(_ codec.Codec, incomingMsg sdk.Msg) []string {
	switch msg := incomingMsg.(type) {
	case *stakingtypes.MsgCreateValidator:
		return []string{msg.ValidatorAddress, msg.DelegatorAddress}
	case *stakingtypes.MsgEditValidator:
		return []string{msg.ValidatorAddress}
	case *stakingtypes.MsgDelegate:
		return []string{msg.DelegatorAddress, msg.ValidatorAddress}
	case *stakingtypes.MsgBeginRedelegate:
		return []string{msg.DelegatorAddress, msg.ValidatorSrcAddress, msg.ValidatorDstAddress}
	case *stakingtypes.MsgUndelegate:
		return []string{msg.DelegatorAddress, msg.ValidatorAddress}
	}

	return nil
}

// IBCTransferMessagesParser returns the list of all the accounts involved in the given
// message if it's related to the x/IBCTransfer module
func IBCTransferMessagesParser(_ codec.Codec, incomingMsg sdk.Msg) []string {
	// nolint:gocritic
	switch msg := incomingMsg.(type) {
	case *ibctransfertypes.MsgTransfer:
		return []string{msg.Sender, msg.Receiver}
	}

	return nil
}

// FeeGrantMessagesParser returns the list of all the accounts involved in the given
// message if it's related to the x/feegrant module
func FeeGrantMessagesParser(_ codec.Codec, incomingMsg sdk.Msg) []string {
	switch msg := incomingMsg.(type) {
	case *feegranttypes.MsgGrantAllowance:
		return []string{msg.Granter, msg.Grantee}
	case *feegranttypes.MsgRevokeAllowance:
		return []string{msg.Granter, msg.Grantee}
	}

	return nil
}

// AuthzMessagesParser returns the list of all the accounts involved in the given
// message if it's related to the x/authz module
func AuthzMessagesParser(_ codec.Codec, incomingMsg sdk.Msg) []string {
	switch msg := incomingMsg.(type) {
	case *authztypes.MsgGrant:
		return []string{msg.Grantee, msg.Granter}
	case *authztypes.MsgRevoke:
		return []string{msg.Grantee, msg.Granter}
	case *authztypes.MsgExec:
		return []string{msg.Grantee}
	}

	return nil
}

// GraphMessagesParser returns the list of all the accounts involved in the given
// message if it's related to the x/graph module
func GraphMessagesParser(_ codec.Codec, incomingMsg sdk.Msg) []string {
	switch msg := incomingMsg.(type) { //nolint:gocritic
	case *graph.MsgCyberlink:
		resp := make([]string, 0, len(msg.Links)*2)
		for _, link := range msg.Links {
			resp = append(resp, link.From, link.To)
		}

		return resp
	}

	return nil
}
