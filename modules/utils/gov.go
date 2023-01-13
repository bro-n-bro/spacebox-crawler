package utils

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
)

func GetProposalContentBytes(content govtypes.Content, cdc codec.Codec) ([]byte, error) {
	// Encode the content properly
	protoContent, ok := content.(proto.Message)
	if !ok {
		return nil, errors.New(fmt.Sprintf("invalid proposal content types: %T", content))
	}

	anyContent, err := codectypes.NewAnyWithValue(protoContent)
	if err != nil {
		return nil, err
	}

	return cdc.MarshalJSON(anyContent)
}
