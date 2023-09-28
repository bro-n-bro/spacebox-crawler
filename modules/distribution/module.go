package distribution

import (
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/rs/zerolog"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	"github.com/bro-n-bro/spacebox-crawler/client/rpc"
	tb "github.com/bro-n-bro/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/bro-n-bro/spacebox-crawler/types"
)

const (
	moduleName = "distribution"
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
	rpcCli *rpc.Client
	broker broker
	tbM    tb.ToBroker
	cdc    codec.Codec
}

func New(b broker, cli *grpcClient.Client, rpcCli *rpc.Client, tbM tb.ToBroker, cdc codec.Codec) *Module {
	l := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("module", moduleName).Logger()

	return &Module{
		log:    &l,
		broker: b,
		client: cli,
		tbM:    tbM,
		cdc:    cdc,
		rpcCli: rpcCli,
	}
}

func (m *Module) Name() string { return moduleName }
