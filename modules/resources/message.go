package resources

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	resources "github.com/cybercongress/go-cyber/x/resources/types"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

const (
	msgErrorPublishingInvestmintMessage = "error while publishing investmint_message message"
)

func (m *Module) HandleMessage(ctx context.Context, index int, bostromMsg sdk.Msg, tx *types.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch msg := bostromMsg.(type) { //nolint:gocritic
	case *resources.MsgInvestmint:
		if err := m.broker.PublishInvestmintMessage(ctx, model.InvestmintMessage{
			Amount:   m.tbM.MapCoin(types.NewCoinFromSDK(msg.Amount)),
			Neuron:   msg.Neuron,
			Resource: msg.Resource,
			Length:   msg.Length,
			TxHash:   tx.TxHash,
			Height:   tx.Height,
			MsgIndex: int64(index),
		}); err != nil {
			return errors.Wrap(err, msgErrorPublishingInvestmintMessage)
		}
	}

	return nil
}
