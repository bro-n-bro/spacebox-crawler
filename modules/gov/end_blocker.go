package bank

import (
	"context"
	"encoding/base64"
	"strconv"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/bro-n-bro/spacebox-crawler/modules/utils"
	"github.com/bro-n-bro/spacebox-crawler/types"
)

var (
	base64KeyProposalID = base64.StdEncoding.EncodeToString([]byte(govtypes.AttributeKeyProposalID))
)

func (m *Module) HandleEndBlocker(ctx context.Context, eventsMap types.BlockerEvents, height int64) error {
	events, ok := eventsMap[govtypes.EventTypeActiveProposal]
	if !ok {
		return nil
	}

	for _, event := range events {
		if len(event.Attributes) < 1 {
			m.log.Warn().
				Int64("height", height).
				Str("handler", "HandleEndBlocker").
				Msg("not enough attributes in event")
			continue
		}

		for _, attr := range event.Attributes {
			if attr.Key == govtypes.AttributeKeyProposalID && attr.Key != base64KeyProposalID {
				continue
			}

			// try to decode value if needed
			if attr.Key == base64KeyProposalID {
				var err error
				attr.Value, err = utils.DecodeToString(attr.Value)
				if err != nil {
					return err
				}
			}

			pID, err := strconv.ParseUint(attr.Value, 10, 64)
			if err != nil {
				m.log.Error().Err(err).Str("handler", "HandleEndBlocker").Msg("parse uint error")
				return err
			}

			if err = m.getAndPublishProposal(ctx, pID, ""); err != nil {
				m.log.Error().
					Err(err).
					Int64("height", height).
					Str("handler", "HandleEndBlocker").
					Msg("get and publish proposal error")
				return err
			}

			return m.getAndPublishTallyResult(ctx, pID, height)
		}
	}

	return nil
}
