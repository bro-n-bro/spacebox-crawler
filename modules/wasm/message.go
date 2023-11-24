package wasm

import (
	"context"
	"slices"

	wasm "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	graph "github.com/cybercongress/go-cyber/x/graph/types"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

const (
	msgErrorPublishingCyberLink = "error while publishing cyber_link"
	msgErrorPublishingParticle  = "error while publishing particle"
)

type clMessage struct {
	CyberLink graph.MsgCyberlink `json:"cyberlink"`
}

func (m *Module) HandleMessage(ctx context.Context, index int, bostromMsg sdk.Msg, tx *types.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	msg, ok := bostromMsg.(*wasm.MsgExecuteContract)
	if !ok {
		return nil
	}

	// try to find links if we were able to unmarshal the message
	clMsg := clMessage{}
	if err := jsoniter.Unmarshal(msg.Msg, &clMsg); err == nil {
		if clMsg.CyberLink.Neuron == "" {
			clMsg.CyberLink.Neuron = msg.Contract
		}

		if clMsg.CyberLink.Neuron != "" && len(clMsg.CyberLink.Links) > 0 {
			for _, link := range clMsg.CyberLink.Links {
				if err = m.broker.PublishCyberlink(ctx, model.Cyberlink{
					ParticleFrom: link.From,
					ParticleTo:   link.To,
					Neuron:       clMsg.CyberLink.Neuron,
					Timestamp:    tx.Timestamp,
					TxHash:       tx.TxHash,
					Height:       tx.Height,
				}); err != nil {
					return errors.Wrap(err, msgErrorPublishingCyberLink)
				}

				for _, particle := range []string{link.From, link.To} {
					if err = m.broker.PublishParticle(ctx, model.Particle{
						Particle:  particle,
						Neuron:    clMsg.CyberLink.Neuron,
						Timestamp: tx.Timestamp,
						TxHash:    tx.TxHash,
						Height:    tx.Height,
					}); err != nil {
						return errors.Wrap(err, msgErrorPublishingParticle)
					}
				}
			}

			return nil
		}
	}

	if err := m.findAndPublishCyberLinks(ctx, tx, index); err != nil {
		if errors.Is(err, types.ErrNoEventFound) || errors.Is(err, types.ErrNoAttributeFound) {
			return nil
		}

		return err
	}

	return nil
}

func (m *Module) findAndPublishCyberLinks(ctx context.Context, tx *types.Tx, index int) error {
	event, err := tx.FindEventByType(index, graph.EventTypeCyberlink)
	if err != nil {
		return err
	}

	links := make([]model.Cyberlink, 0, 1)

	switch {
	case len(event.Attributes) <= 3:
		links = append(links, m.findOneLink(event, tx))
	case event.Attributes[2].Key == graph.AttributeKeyNeuron:
		slices.Grow(links, len(event.Attributes)/3)
		links = append(links, m.findWithSequence(event, tx)...)
	case event.Attributes[len(event.Attributes)-1].Key == graph.AttributeKeyNeuron:
		slices.Grow(links, len(event.Attributes)/2)
		links = append(links, m.findWithCommonNeuron(event, tx)...)
	}

	for _, link := range links {
		if err = m.broker.PublishCyberlink(ctx, link); err != nil {
			return errors.Wrap(err, msgErrorPublishingCyberLink)
		}

		for _, particle := range []string{link.ParticleFrom, link.ParticleTo} {
			if err = m.broker.PublishParticle(ctx, model.Particle{
				Particle:  particle,
				Neuron:    link.Neuron,
				Timestamp: tx.Timestamp,
				TxHash:    tx.TxHash,
				Height:    tx.Height,
			}); err != nil {
				return errors.Wrap(err, msgErrorPublishingParticle)
			}
		}
	}

	return nil
}

func (m *Module) findOneLink(event sdk.StringEvent, tx *types.Tx) model.Cyberlink {
	from, _ := tx.FindAttributeByKey(event, graph.AttributeKeyParticleFrom)
	to, _ := tx.FindAttributeByKey(event, graph.AttributeKeyParticleTo)
	neuron, _ := tx.FindAttributeByKey(event, graph.AttributeKeyNeuron)

	return model.Cyberlink{
		ParticleFrom: from,
		ParticleTo:   to,
		Neuron:       neuron,
		Timestamp:    tx.Timestamp,
		TxHash:       tx.TxHash,
		Height:       tx.Height,
	}
}

//	{
//		key: "particleFrom",
//		value: "QmR8xA9EyCQhGWu9cUs4fmsZnsVGGFLfY8z1eHrbqNRQUB"
//	},
//
//	{
//		 key: "particleTo",
//		value: "QmW5GREog52duzpQbHQ2da8NCQSkL218TWW9hQzsR11bGM"
//	},
//
//	{
//		key: "particleFrom",
//		value: "QmaoRBnpjnjcqfYyUDvQx2s7ZBtHJbycU8bEw9AWWLaQVd"
//	},
//
//	{
//		key: "particleTo",
//		value: "QmW5GREog52duzpQbHQ2da8NCQSkL218TWW9hQzsR11bGM"
//	},
//
//	{
//		key: "neuron",
//		value: "bostrom1jkte0pytr85qg0whmgux3vmz9ehmh82w40h8gaqeg435fnkyfxqq9qaku3"
//	}
func (m *Module) findWithCommonNeuron(event sdk.StringEvent, tx *types.Tx) []model.Cyberlink {
	links := make([]model.Cyberlink, 0)

	// common neuron
	neuron := event.Attributes[len(event.Attributes)-1].Value

	for i := 0; i < len(event.Attributes)-1; i += 2 {
		var from, to string
		switch event.Attributes[i].Key {
		case graph.AttributeKeyParticleFrom:
			from = event.Attributes[i].Value
		case graph.AttributeKeyParticleTo:
			to = event.Attributes[i].Value
		}

		switch event.Attributes[i+1].Key {
		case graph.AttributeKeyParticleFrom:
			from = event.Attributes[i+1].Value
		case graph.AttributeKeyParticleTo:
			to = event.Attributes[i+1].Value
		}

		links = append(links, model.Cyberlink{
			ParticleFrom: from,
			ParticleTo:   to,
			Neuron:       neuron,
			Timestamp:    tx.Timestamp,
			TxHash:       tx.TxHash,
			Height:       tx.Height,
		})
	}

	return links
}

//	{
//		key: "particleFrom",
//		value: "QmTEv3kLMPNkX3yv92ixzr17t5ayVi99Nkso3QDtKs7qta"
//	},
//
//	{
//		key: "particleTo",
//		value: "QmQ8XsZSdNWNbu1FZpJyua8CXMQTq6dLnsHaqLtpEw8GXL"
//	},
//
//	{
//		key: "neuron",
//		value: "bostrom1jkte0pytr85qg0whmgux3vmz9ehmh82w40h8gaqeg435fnkyfxqq9qaku3"
//	},
//
//	{
//		key: "particleFrom",
//		value: "QmTEv3kLMPNkX3yv92ixzr17t5ayVi99Nkso3QDtKs7qta"
//	},
//
//	{
//		key: "particleTo",
//		value: "QmQ8XsZSdNWNbu1FZpJyua8CXMQTq6dLnsHaqLtpEw8GXL"
//	},
//
//	{
//		key: "neuron",
//		value: "bostrom1jkte0pytr85qg0whmgux3vmz9ehmh82w40h8gaqeg435fnkyfxqq9qaku3"
//	}
func (m *Module) findWithSequence(event sdk.StringEvent, tx *types.Tx) []model.Cyberlink {
	links := make([]model.Cyberlink, 0)

	for i := 0; i < len(event.Attributes); i += 3 {
		var from, to, neuron string
		switch event.Attributes[i].Key {
		case graph.AttributeKeyParticleFrom:
			from = event.Attributes[i].Value
		case graph.AttributeKeyParticleTo:
			to = event.Attributes[i].Value
		case graph.AttributeKeyNeuron:
			neuron = event.Attributes[i].Value
		}

		switch event.Attributes[i+1].Key {
		case graph.AttributeKeyParticleFrom:
			from = event.Attributes[i+1].Value
		case graph.AttributeKeyParticleTo:
			to = event.Attributes[i+1].Value
		case graph.AttributeKeyNeuron:
			neuron = event.Attributes[i+1].Value
		}

		switch event.Attributes[i+2].Key {
		case graph.AttributeKeyParticleFrom:
			from = event.Attributes[i+2].Value
		case graph.AttributeKeyParticleTo:
			to = event.Attributes[i+2].Value
		case graph.AttributeKeyNeuron:
			neuron = event.Attributes[i+2].Value
		}

		links = append(links, model.Cyberlink{
			ParticleFrom: from,
			ParticleTo:   to,
			Neuron:       neuron,
			Timestamp:    tx.Timestamp,
			TxHash:       tx.TxHash,
			Height:       tx.Height,
		})
	}

	return links
}
