package graph

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	graph "github.com/cybercongress/go-cyber/x/graph/types"

	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

func (m *Module) HandleMessage(ctx context.Context, index int, cosmosMsg sdk.Msg, tx *types.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch msg := cosmosMsg.(type) { //nolint:gocritic
	case *graph.MsgCyberlink:
		for i, link := range msg.Links {
			if err := m.broker.PublishCyberLinkMessage(ctx, model.CyberLinkMessage{
				ParticleFrom: link.From,
				ParticleTo:   link.To,
				Neuron:       msg.Neuron,
				TxHash:       tx.TxHash,
				Height:       tx.Height,
				MsgIndex:     int64(index),
				LinkIndex:    int64(i), // TODO: check if it is correct
			}); err != nil {
				m.log.Err(err).Int64("height", tx.Height).Msg("error while publishing cyberlink message")
				return err
			}
		}
		// TODO: add other messages
	}

	return nil
}
