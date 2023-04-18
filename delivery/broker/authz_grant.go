package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishAuthzGrant(ctx context.Context, grant model.AuthzGrant) error {
	data, err := jsoniter.Marshal(grant)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(AuthzGrant, data)
}
