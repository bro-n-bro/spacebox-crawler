package ibc

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

type broker interface {
	PublishTransferMessage(context.Context, model.TransferMessage) error
	PublishAcknowledgementMessage(context.Context, model.AcknowledgementMessage) error
	PublishReceivePacketMessage(context.Context, model.RecvPacketMessage) error
	PublishDenomTrace(context.Context, model.DenomTrace) error
}
