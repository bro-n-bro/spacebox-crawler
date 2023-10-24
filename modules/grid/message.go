package grid

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	grid "github.com/cybercongress/go-cyber/x/grid/types"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

		return m.getAndPublishRoute(ctx, tx.TxHash, msg.Source, msg.Destination, tx.Timestamp, tx.Height)

	case *grid.MsgEditRoute:
		if err := m.broker.PublishEditRouteMessage(ctx, model.EditRouteMessage{
			Value:       m.tbM.MapCoin(types.NewCoinFromSDK(msg.Value)),
			Source:      msg.Source,
			Destination: msg.Destination,
			TxHash:      tx.TxHash,
			Height:      tx.Height,
			MsgIndex:    int64(index),
		}); err != nil {
			return errors.Wrap(err, msgErrorPublishingEditRouteMessage)
		}

		return m.getAndPublishRoute(ctx, tx.TxHash, msg.Source, msg.Destination, tx.Timestamp, tx.Height)

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

		return m.getAndPublishRoute(ctx, tx.TxHash, msg.Source, msg.Destination, tx.Timestamp, tx.Height)

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

		return m.getAndPublishRoute(ctx, tx.TxHash, msg.Source, msg.Destination, tx.Timestamp, tx.Height)
	}

	return nil
}

func (m *Module) getAndPublishRoute(ctx context.Context, tx, source, destination, ts string, height int64) error {
	var (
		key      = source + "-" + destination
		isActive = true
		value    model.Coins
		alias    string
	)

	// publish only newest heights
	if m.routeCache != nil && !m.routeCache.UpdateCacheValue(key, height) {
		return nil
	}

	route, err := m.client.GridQueryClient.Route(ctx, &grid.QueryRouteRequest{
		Source:      source,
		Destination: destination,
	})

	if err != nil {
		if s, ok := status.FromError(err); !ok || s.Code() != codes.NotFound {
			m.log.Error().Err(err).
				Str("source", source).
				Str("destination", destination).
				Msg("error while getting route")

			return err
		}

		isActive = false
	}

	// case when route is not found
	if route != nil {
		value = m.tbM.MapCoins(types.NewCoinsFromSDK(route.Route.Value))
		alias = route.Route.Name
	}

	return m.broker.PublishRoute(ctx, model.Route{
		Value:       value,
		Source:      source,
		Destination: destination,
		Alias:       alias,
		Timestamp:   ts,
		TxHash:      tx,
		Height:      height,
		IsActive:    isActive,
	})
}
