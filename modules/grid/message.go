package grid

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	grid "github.com/cybercongress/go-cyber/x/grid/types"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

const (
	msgErrorPublishingCreateRouteMessage   = "error while publishing create_route message"
	msgErrorPublishingEditRouteMessage     = "error while publishing edit_route message"
	msgErrorPublishingEditRouteNameMessage = "error while publishing edit_route_name message"
	msgErrorPublishingDeleteRouteMessage   = "error while publishing delete_route message"
)

func (m *Module) HandleMessage(ctx context.Context, index int, bostromMsg sdk.Msg, tx *types.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch msg := bostromMsg.(type) {
	case *grid.MsgCreateRoute:
		if err := m.broker.PublishCreateRouteMessage(ctx, model.CreateRouteMessage{
			Source:      msg.Source,
			Destination: msg.Destination,
			Name:        msg.Name,
			TxHash:      tx.TxHash,
			Height:      tx.Height,
			MsgIndex:    int64(index),
		}); err != nil {
			return errors.Wrap(err, msgErrorPublishingCreateRouteMessage)
		}

	case *grid.MsgEditRoute:
		if err := m.broker.PublishEditRouteMessage(ctx, model.EditRouteMessage{
			Value:       m.tbM.MapCoin(types.NewCoinFromCdk(msg.Value)),
			Source:      msg.Source,
			Destination: msg.Destination,
			TxHash:      tx.TxHash,
			Height:      tx.Height,
			MsgIndex:    int64(index),
		}); err != nil {
			return errors.Wrap(err, msgErrorPublishingEditRouteMessage)
		}

	case *grid.MsgEditRouteName:
		if err := m.broker.PublishEditRouteNameMessage(ctx, model.EditRouteNameMessage{
			Source:      msg.Source,
			Destination: msg.Destination,
			Name:        msg.Name,
			TxHash:      tx.TxHash,
			Height:      tx.Height,
			MsgIndex:    int64(index),
		}); err != nil {
			return errors.Wrap(err, msgErrorPublishingEditRouteNameMessage)
		}

	case *grid.MsgDeleteRoute:
		if err := m.broker.PublishDeleteRouteMessage(ctx, model.DeleteRouteMessage{
			Source:      msg.Source,
			Destination: msg.Destination,
			TxHash:      tx.TxHash,
			Height:      tx.Height,
			MsgIndex:    int64(index),
		}); err != nil {
			return errors.Wrap(err, msgErrorPublishingDeleteRouteMessage)
		}
	}

	return nil
}
