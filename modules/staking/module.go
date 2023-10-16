package staking

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/rs/zerolog"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	"github.com/bro-n-bro/spacebox-crawler/modules/utils"
	tb "github.com/bro-n-bro/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/bro-n-bro/spacebox-crawler/types"
)

const (
	ModuleName = "staking"
)

var (
	_ types.Module            = &Module{}
	_ types.GenesisHandler    = &Module{}
	_ types.BlockHandler      = &Module{}
	_ types.MessageHandler    = &Module{}
	_ types.ValidatorsHandler = &Module{}
)

type (
	AccountCache[K, V comparable] interface {
		UpdateCacheValue(K, V) bool
	}

	Module struct {
		log            *zerolog.Logger
		client         *grpcClient.Client
		broker         broker
		tbM            tb.ToBroker
		cdc            codec.Codec
		accCache       AccountCache[string, int64]
		enabledModules []string // xxx fixme
	}
)

func New(b broker, cli *grpcClient.Client, tbM tb.ToBroker, cdc codec.Codec, modules []string) *Module {
	return &Module{
		log:            utils.NewModuleLogger(ModuleName),
		broker:         b,
		client:         cli,
		tbM:            tbM,
		cdc:            cdc,
		enabledModules: modules,
	}
}

func (m *Module) Name() string { return ModuleName }

func (m *Module) WithAccountCache(cache AccountCache[string, int64]) *Module {
	if cache != nil {
		m.accCache = cache
	}

	return m
}
