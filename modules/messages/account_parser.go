package messages

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v5/modules/apps/transfer/types"
)

// 07-tendermint-259
// ibc.lightclients.tendermint.v1.Header

// MessageNotSupported returns an error telling that the given message is not supported
func MessageNotSupported(msg sdk.Msg) error {
	defer func() { // FIXME
		if r := recover(); r != nil {
			println("message type not supported:", r)
		}
	}()
	// return fmt.Errorf("message type not supported: %s", msg.String())
	return nil
}

// MessageAddressesParser represents a function that extracts all the
// involved addresses from a provided message (both accounts and validators)
type MessageAddressesParser = func(cdc codec.Codec, msg sdk.Msg) ([]string, error)

// JoinMessageParsers joins together all the given parsers, calling them in order
func JoinMessageParsers(parsers ...MessageAddressesParser) MessageAddressesParser {
	return func(cdc codec.Codec, msg sdk.Msg) ([]string, error) {
		for _, parser := range parsers {
			// Try getting the addresses
			addresses, _ := parser(cdc, msg)

			// If some addresses are found, return them
			if len(addresses) > 0 {
				return addresses, nil
			}
		}
		return nil, MessageNotSupported(msg)
	}
}

// CosmosMessageAddressesParser represents a MessageAddressesParser that parses a
// Cosmos message and returns all the involved addresses (both accounts and validators)
var CosmosMessageAddressesParser = JoinMessageParsers(
	BankMessagesParser,
	CrisisMessagesParser,
	DistributionMessagesParser,
	EvidenceMessagesParser,
	GovMessagesParser,
	IBCTransferMessagesParser,
	SlashingMessagesParser,
	StakingMessagesParser,
	DefaultMessagesParser,
)

// DefaultMessagesParser represents the default messages parser that simply returns the list
// of all the signers of a message
func DefaultMessagesParser(_ codec.Codec, cosmosMsg sdk.Msg) ([]string, error) {
	cosmosSigners := cosmosMsg.GetSigners()
	signers := make([]string, len(cosmosSigners))
	for index, signer := range cosmosSigners {
		signers[index] = signer.String()
	}
	return signers, nil
}

// BankMessagesParser returns the list of all the accounts involved in the given
// message if it's related to the x/bank module
func BankMessagesParser(_ codec.Codec, cosmosMsg sdk.Msg) ([]string, error) {
	switch msg := cosmosMsg.(type) {
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

	return nil, MessageNotSupported(cosmosMsg)
}

// CrisisMessagesParser returns the list of all the accounts involved in the given
// message if it's related to the x/crisis module
func CrisisMessagesParser(_ codec.Codec, cosmosMsg sdk.Msg) ([]string, error) {
	// nolint:gocritic
	switch msg := cosmosMsg.(type) {
	case *crisistypes.MsgVerifyInvariant:
		return []string{msg.Sender}, nil
	}

	return nil, MessageNotSupported(cosmosMsg)
}

// DistributionMessagesParser returns the list of all the accounts involved in the given
// message if it's related to the x/distribution module
func DistributionMessagesParser(_ codec.Codec, cosmosMsg sdk.Msg) ([]string, error) {
	switch msg := cosmosMsg.(type) {

	case *distrtypes.MsgSetWithdrawAddress:
		return []string{msg.DelegatorAddress, msg.WithdrawAddress}, nil

	case *distrtypes.MsgWithdrawDelegatorReward:
		return []string{msg.DelegatorAddress, msg.ValidatorAddress}, nil

	case *distrtypes.MsgWithdrawValidatorCommission:
		return []string{msg.ValidatorAddress}, nil

	case *distrtypes.MsgFundCommunityPool:
		return []string{msg.Depositor}, nil

	}

	return nil, MessageNotSupported(cosmosMsg)
}

// EvidenceMessagesParser returns the list of all the accounts involved in the given
// message if it's related to the x/evidence module
func EvidenceMessagesParser(_ codec.Codec, cosmosMsg sdk.Msg) ([]string, error) {
	// nolint:gocritic
	switch msg := cosmosMsg.(type) {
	case *evidencetypes.MsgSubmitEvidence:
		return []string{msg.Submitter}, nil
	}

	return nil, MessageNotSupported(cosmosMsg)
}

// GovMessagesParser returns the list of all the accounts involved in the given
// message if it's related to the x/gov module
func GovMessagesParser(cdc codec.Codec, cosmosMsg sdk.Msg) ([]string, error) {
	switch msg := cosmosMsg.(type) {

	case *govtypes.MsgSubmitProposal:
		addresses := []string{msg.Proposer}

		var content govtypes.Content
		err := cdc.UnpackAny(msg.Content, &content)
		if err != nil {
			return nil, err
		}

		// nolint:gocritic
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

	return nil, MessageNotSupported(cosmosMsg)
}

//
//// IBCTransferMessagesParser returns the list of all the accounts involved in the given
//// message if it's related to the x/iBCTransfer module
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
func SlashingMessagesParser(_ codec.Codec, cosmosMsg sdk.Msg) ([]string, error) {
	// nolint:gocritic
	switch msg := cosmosMsg.(type) {
	case *slashingtypes.MsgUnjail:
		return []string{msg.ValidatorAddr}, nil

	}

	return nil, MessageNotSupported(cosmosMsg)
}

// StakingMessagesParser returns the list of all the accounts involved in the given
// message if it's related to the x/staking module
func StakingMessagesParser(_ codec.Codec, cosmosMsg sdk.Msg) ([]string, error) {
	switch msg := cosmosMsg.(type) {

	case *stakingtypes.MsgCreateValidator:
		return []string{msg.ValidatorAddress, msg.DelegatorAddress}, nil

	case *stakingtypes.MsgEditValidator:
		return []string{msg.ValidatorAddress}, nil

	case *stakingtypes.MsgDelegate:
		return []string{msg.DelegatorAddress, msg.ValidatorAddress}, nil

	case *stakingtypes.MsgBeginRedelegate:
		return []string{msg.DelegatorAddress, msg.ValidatorSrcAddress, msg.ValidatorDstAddress}, nil

	case *stakingtypes.MsgUndelegate:
		return []string{msg.DelegatorAddress, msg.ValidatorAddress}, nil

	}

	return nil, MessageNotSupported(cosmosMsg)
}

// IBCTransferMessagesParser returns the list of all the accounts involved in the given
// message if it's related to the x/iBCTransfer module
func IBCTransferMessagesParser(_ codec.Codec, cosmosMsg sdk.Msg) ([]string, error) {
	// nolint:gocritic
	switch msg := cosmosMsg.(type) {
	case *ibctransfertypes.MsgTransfer:
		return []string{msg.Sender, msg.Receiver}, nil
	}

	return nil, MessageNotSupported(cosmosMsg)
}
