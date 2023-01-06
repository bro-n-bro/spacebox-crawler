package mint

import (
	"os"

	grpcClient "github.com/hexy-dev/spacebox-crawler/client/grpc"
	"github.com/hexy-dev/spacebox-crawler/internal/rep"
	"github.com/hexy-dev/spacebox-crawler/types"

	tb "github.com/hexy-dev/spacebox-crawler/pkg/mapper/to_broker"

	"github.com/rs/zerolog"
)

var (
	_ types.Module       = &Module{}
	_ types.BlockHandler = &Module{}
)

type Module struct {
	log    *zerolog.Logger
	client *grpcClient.Client
	broker broker
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
