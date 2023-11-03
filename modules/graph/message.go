package graph

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	graph "github.com/cybercongress/go-cyber/x/graph/types"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

const (
	msgErrorPublishingCyberLinkMessage = "error while publishing cyber_link message"
	msgErrorPublishingCyberLink        = "error while publishing cyber_link"
	msgErrorPublishingParticle         = "error while publishing particle"
)

func (m *Module) HandleMessage(ctx context.Context, index int, bostromMsg sdk.Msg, tx *types.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch msg := bostromMsg.(type) { //nolint:gocritic
	case *graph.MsgCyberlink:
		for i, link := range msg.Links {
			if err := m.broker.PublishCyberlinkMessage(ctx, model.CyberlinkMessage{
				ParticleFrom: link.From,    //
				ParticleTo:   link.To,      //
				Neuron:       msg.Neuron,   //
				TxHash:       tx.TxHash,    //
				Height:       tx.Height,    //
				MsgIndex:     int64(index), //
				LinkIndex:    int64(i),     //
			}); err != nil {
				return errors.Wrap(err, msgErrorPublishingCyberLinkMessage)
			}

			if err := m.broker.PublishCyberlink(ctx, model.Cyberlink{
				ParticleFrom: link.From,    //
				ParticleTo:   link.To,      //
				Neuron:       msg.Neuron,   //
				TxHash:       tx.TxHash,    //
				Height:       tx.Height,    //
				Timestamp:    tx.Timestamp, //
			}); err != nil {
				return errors.Wrap(err, msgErrorPublishingCyberLink)
			}

			for _, particle := range []string{link.From, link.To} {
				if err := m.broker.PublishParticle(ctx, model.Particle{
					Particle:  particle,     //
					Neuron:    msg.Neuron,   //
					Timestamp: tx.Timestamp, //
					TxHash:    tx.TxHash,    //
					Height:    tx.Height,    //
				}); err != nil {
					return errors.Wrap(err, msgErrorPublishingParticle)
				}
			}
		}
	}

	return nil
}
