package slasing

import (
	"os"

	grpcClient "github.com/hexy-dev/spacebox-crawler/client/grpc"
	"github.com/hexy-dev/spacebox-crawler/internal/rep"
	"github.com/hexy-dev/spacebox-crawler/types"

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
