package broker

import (
	"context"

	banktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/segmentio/kafka-go"
)

func (b *Broker) PublishBank(ctx context.Context, response banktypes.Coins) error {

	jsonBytes, err := response.MarshalJSON()
	if err != nil {
		b.log.Error().Err(err).Msg("json marshal fail")
		return err
	}

	err = b.writer.WriteMessages(ctx,
		// NOTE: Each Message has Topic defined, otherwise an error is returned.
		kafka.Message{
			Topic: "block-module",
			Key:   []byte("block"),
			Value: jsonBytes,
		},
	)
	if err != nil {
		b.log.Error().Err(err).Msg("failed to write messages")
		return err
	}
	return nil
}
