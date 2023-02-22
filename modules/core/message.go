package core

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"

	"github.com/bro-n-bro/spacebox-crawler/types"
)

func (m *Module) HandleMessage(ctx context.Context, index int, msg sdk.Msg, tx *types.Tx) error {
	// Marshal the value properly
	msgValue, err := m.cdc.MarshalJSON(msg)
	if err != nil {
		return err
	}

	return m.broker.PublishMessage(ctx,
		m.tbM.MapMessage(tx.TxHash, proto.MessageName(msg), tx.Signer, index, m.parser(m.cdc, msg), msgValue))
}
