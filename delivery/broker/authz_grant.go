package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishAuthzGrant(ctx context.Context, grant model.AuthzGrant) error {
	return b.marshalAndProduce(AuthzGrant, grant)
}
