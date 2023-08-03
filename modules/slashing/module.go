package slashing

import (
	"os"

	"github.com/rs/zerolog"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	tb "github.com/bro-n-bro/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/bro-n-bro/spacebox-crawler/types"
)

const (
	moduleName = "slashing"
)

var (
	_ types.Module              = &Module{}
	_ types.BlockHandler        = &Module{}
	_ types.MessageHandler      = &Module{}
	_ types.BeginBlockerHandler = &Module{}
)

type Module struct {
	log    *zerolog.Logger
	client *grpcClient.Client
	tbM    tb.ToBroker
	broker broker
}

func New(b broker, cli *grpcClient.Client, tbM tb.ToBroker) *Module {
	l := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("module", moduleName).Logger()

	return &Module{
		log:    &l,
		client: cli,
		tbM:    tbM,
		broker: b,
	}
}

func (m *Module) Name() string { return moduleName }
