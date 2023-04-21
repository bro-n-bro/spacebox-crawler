package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishVoteWeightedMessage(_ context.Context, vwm model.VoteWeightedMessage) error {
	data, err := jsoniter.Marshal(vwm)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(VoteWeightedMessage, data)
}
