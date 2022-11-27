package slasing

import (
	grpcClient "bro-n-bro-osmosis/client/grpc"
	"bro-n-bro-osmosis/internal/rep"
	"bro-n-bro-osmosis/types"
	"os"

	"github.com/rs/zerolog"
)

var (
	_ types.Module      = &Module{}
	_ types.BlockModule = &Module{}
)

type Module struct {
	log    *zerolog.Logger
	broker rep.Broker
	client *grpcClient.Client
}

func New(b rep.Broker, cli *grpcClient.Client) *Module {
	l := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("module", "slashing").Logger()

	return &Module{
		log:    &l,
		broker: b,
		client: cli,
	}
}

func (m *Module) Name() string { return "slashing" }
