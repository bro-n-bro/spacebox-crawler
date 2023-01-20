package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishCommunityPool(ctx context.Context, cp model.CommunityPool) error {
	data, err := jsoniter.Marshal(cp)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(CommunityPool, data)
}
