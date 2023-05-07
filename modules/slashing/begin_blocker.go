package slashing

import (
	"context"

	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"

	"github.com/bro-n-bro/spacebox-crawler/types"
)

func (m *Module) HandleBeginBlocker(ctx context.Context, eventsMap types.BlockerEvents, height int64) error {
	events, ok := eventsMap[slashingtypes.EventTypeSlash]
	if !ok {
		return nil
	}

	for _, e := range events {
		_ = e
	}

	return nil
}
