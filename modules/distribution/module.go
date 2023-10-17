package distribution

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/rs/zerolog"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	"github.com/bro-n-bro/spacebox-crawler/client/rpc"
	"github.com/bro-n-bro/spacebox-crawler/modules/utils"
	tb "github.com/bro-n-bro/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/bro-n-bro/spacebox-crawler/types"
)

const (
	ModuleName = "distribution"
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
	return &Module{
		log:    utils.NewModuleLogger(ModuleName),
		broker: b,
		client: cli,
		tbM:    tbM,
		cdc:    cdc,
		rpcCli: rpcCli,
	}
}

func (m *Module) Name() string { return ModuleName }
