package authz

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/rs/zerolog"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	"github.com/bro-n-bro/spacebox-crawler/modules/utils"
	tb "github.com/bro-n-bro/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/bro-n-bro/spacebox-crawler/types"
)

const (
	ModuleName = "authz"
)

var (
	_ types.Module         = &Module{}
	_ types.MessageHandler = &Module{}
)

// Module is a module for authz.
type Module struct {
	log    *zerolog.Logger
	client *grpcClient.Client
	broker broker
	tbM    tb.ToBroker
	cdc    codec.Codec
}

// New creates a new authz module.
func New(b broker, cli *grpcClient.Client, tb tb.ToBroker, cdc codec.Codec) *Module {
	return &Module{
		log:    utils.NewModuleLogger(ModuleName),
		broker: b,
		tbM:    tb,
		client: cli,
		cdc:    cdc,
	}
}

// Name returns the module name.
func (m *Module) Name() string { return ModuleName }
