package wasm

import (
	"context"

	wasm "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	graph "github.com/cybercongress/go-cyber/x/graph/types"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

const (
	msgErrorPublishingCyberLink = "error while publishing cyber_link"
	msgErrorPublishingParticle  = "error while publishing particle"
)

func (m *Module) HandleMessage(ctx context.Context, index int, bostromMsg sdk.Msg, tx *types.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	if _, ok := bostromMsg.(*wasm.MsgExecuteContract); !ok {
		return nil
	}

	if err := m.findAndPublishCyberLink(ctx, tx, index); err != nil {
		if errors.Is(err, types.ErrNoEventFound) || errors.Is(err, types.ErrNoAttributeFound) {
			return nil
		}

		return err
	}

	return nil
}

func (m *Module) findAndPublishCyberLink(ctx context.Context, tx *types.Tx, index int) error {
	event, err := tx.FindEventByType(index, graph.EventTypeCyberlink)
	if err != nil {
		return err
	}

	from, err := tx.FindAttributeByKey(event, graph.AttributeKeyParticleFrom)
	if err != nil {
		return err
	}

	to, err := tx.FindAttributeByKey(event, graph.AttributeKeyParticleTo)
	if err != nil {
		return err
	}

	neuron, err := tx.FindAttributeByKey(event, graph.AttributeKeyNeuron)
	if err != nil {
		return err
	}

	if err = m.broker.PublishCyberlink(ctx, model.Cyberlink{
		ParticleFrom: from,
		ParticleTo:   to,
		Neuron:       neuron,
		Timestamp:    tx.Timestamp,
		TxHash:       tx.TxHash,
		Height:       tx.Height,
	}); err != nil {
		return errors.Wrap(err, msgErrorPublishingCyberLink)
	}

	for _, particle := range []string{from, to} {
		if err = m.broker.PublishParticle(ctx, model.Particle{
			Particle:  particle,
			Neuron:    neuron,
			Timestamp: tx.Timestamp,
			TxHash:    tx.TxHash,
			Height:    tx.Height,
		}); err != nil {
			return errors.Wrap(err, msgErrorPublishingParticle)
		}
	}

	return nil
}
