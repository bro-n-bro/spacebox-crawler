package modules

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AddressesParser represents a function that extracts all the
// involved addresses from a provided message (both accounts and validators)
type AddressesParser = func(cdc codec.Codec, msg sdk.Msg) ([]string, error)

// JoinMessageParsers joins together all the given parsers, calling them in order
func JoinMessageParsers(parsers ...AddressesParser) AddressesParser {
	return func(cdc codec.Codec, msg sdk.Msg) ([]string, error) {
		for _, parser := range parsers {
			// Try getting the addresses and return them
			if addresses, _ := parser(cdc, msg); len(addresses) > 0 {
				return addresses, nil
			}
		}
		return nil, nil
	}
}

// DefaultMessagesParser represents the default messages parser that simply returns the list
// of all the signers of a message
func DefaultMessagesParser(_ codec.Codec, cosmosMsg sdk.Msg) ([]string, error) {
	var signers = make([]string, len(cosmosMsg.GetSigners()))
	for index, signer := range cosmosMsg.GetSigners() {
		signers[index] = signer.String()
	}
	return signers, nil
}
