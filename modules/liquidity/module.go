package liquidity

import (
	"os"

	"github.com/rs/zerolog"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	tb "github.com/bro-n-bro/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/bro-n-bro/spacebox-crawler/types"
)

const (
	moduleName = "liquidity"
)

var (
	_ types.Module            = &Module{}
	_ types.EndBlockerHandler = &Module{}
)

type (
	Module struct {
		log    *zerolog.Logger
		client *grpcClient.Client
		broker broker
		tbM    tb.ToBroker
	}
)

func New(b broker, cli *grpcClient.Client, tbM tb.ToBroker) *Module {
	l := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("module", moduleName).Logger()

	m := &Module{
		log:    &l,
		broker: b,
		client: cli,
		tbM:    tbM,
	}

	return m
}

func (m *Module) Name() string { return moduleName }
