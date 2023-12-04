package bank

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/rs/zerolog"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	"github.com/bro-n-bro/spacebox-crawler/modules/utils"
	tb "github.com/bro-n-bro/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/bro-n-bro/spacebox-crawler/types"
)

const (
	ModuleName = "gov"
)

var (
	_ types.Module            = &Module{}
	_ types.GenesisHandler    = &Module{}
	_ types.BlockHandler      = &Module{}
	_ types.MessageHandler    = &Module{}
	_ types.EndBlockerHandler = &Module{}
)

type (
	TallyCache[K, V comparable] interface {
		UpdateCacheValue(K, V) bool
	}

	Module struct {
		log    *zerolog.Logger
		client *grpcClient.Client
		broker broker
		tbM    tb.ToBroker
		cdc    codec.Codec

		tallyCache TallyCache[uint64, int64]
	}
)

func New(b broker, cli *grpcClient.Client, tbM tb.ToBroker, cdc codec.Codec) *Module {
	m := &Module{
		log:    utils.NewModuleLogger(ModuleName),
		broker: b,
		client: cli,
		tbM:    tbM,
		cdc:    cdc,
	}

	return m
}

func (m *Module) Name() string { return ModuleName }

func (m *Module) WithTallyCache(cache TallyCache[uint64, int64]) *Module {
	if cache != nil {
		m.tallyCache = cache
	}

	return m
}
