package core

import (
	"context"

	"github.com/gogo/protobuf/proto"

	"github.com/hexy-dev/spacebox-crawler/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m *Module) HandleMessage(ctx context.Context, index int, msg sdk.Msg, tx *types.Tx) error {
	// Get the involved addresses
	addresses, err := m.parser(m.cdc, msg)
	if err != nil {
		return err
	}

	// Marshal the value properly
	msgValue, err := m.cdc.MarshalJSON(msg)
	if err != nil {
		return err
	}

	// msg.GetSigners() TODO:
	return m.broker.PublishMessage(ctx, m.tbM.MapMessage(tx.TxHash, proto.MessageName(msg), tx.Signer, index, addresses, msgValue))
}
