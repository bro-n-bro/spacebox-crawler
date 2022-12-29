package messages

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/hexy-dev/spacebox-crawler/types"
)

func (m *Module) HandleMessage(_ context.Context, _ int, cdkMsg sdk.Msg, _ *types.Tx) error {
	// Get the involved addresses
	_, err := m.parser(m.cdc, cdkMsg)
	if err != nil {
		return err
	}

	// Marshal the value properly
	_, err = m.cdc.MarshalJSON(cdkMsg)
	if err != nil {
		return err
	}

	// return db.SaveMessage(types.NewMessage(
	//	tx.TxHash,
	//	index,
	//	proto.MessageName(msg),
	//	string(bz),
	//	addresses,
	// ))
	return nil
}
