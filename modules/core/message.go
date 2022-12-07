package core

import (
	"context"

	"github.com/gogo/protobuf/proto"

	"bro-n-bro-osmosis/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m *Module) HandleMessage(ctx context.Context, index int, msg sdk.Msg, tx *types.Tx) error {
	// Get the involved addresses
	addresses, err := m.parser(m.cdc, msg)
	if err != nil {
		return err
	}

	// Marshal the value properly
	bz, err := m.cdc.MarshalJSON(msg)
	if err != nil {
		return err
	}

	// msg.GetSigners() TODO:

	err = m.broker.PublishMessage(ctx,
		m.tbM.MapMessage(tx.TxHash, proto.MessageName(msg), "", index, addresses, bz))

	return err
}
