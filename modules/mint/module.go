package mint

import (
	grpcClient "bro-n-bro-osmosis/client/grpc"
	"bro-n-bro-osmosis/internal/rep"
	"bro-n-bro-osmosis/types"
	"os"

	tb "bro-n-bro-osmosis/pkg/mapper/to_broker"

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
	tbM    tb.ToBroker
}

func New(b rep.Broker, cli *grpcClient.Client, tbM tb.ToBroker) *Module {
	l := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("module", "mint").Logger()

	return &Module{
		log:    &l,
		broker: b,
		client: cli,
		tbM:    tbM,
	}
}

func (m *Module) Name() string { return "mint" }
