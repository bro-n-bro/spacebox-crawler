package auth

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/rs/zerolog"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	"github.com/bro-n-bro/spacebox-crawler/modules/core"
	"github.com/bro-n-bro/spacebox-crawler/modules/utils"
	tb "github.com/bro-n-bro/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/bro-n-bro/spacebox-crawler/types"
)

const (
	ModuleName = "auth"
)

var (
	_ types.Module         = &Module{}
	_ types.GenesisHandler = &Module{}
	_ types.MessageHandler = &Module{}
)

type (
	AccountCache[K, V comparable] interface {
		UpdateCacheValue(K, V) bool
	}

	Module struct {
		log      *zerolog.Logger
		client   *grpcClient.Client
		broker   broker
		tbM      tb.ToBroker
		cdc      codec.Codec
		parser   core.MsgAddrParser
		accCache AccountCache[string, int64]
	}
)

func New(b broker, cli *grpcClient.Client, tb tb.ToBroker, cdc codec.Codec, parser core.MsgAddrParser) *Module {
	return &Module{
		log:    utils.NewModuleLogger(ModuleName),
		broker: b,
		tbM:    tb,
		client: cli,
		cdc:    cdc,
		parser: parser,
	}
}

func (m *Module) Name() string { return ModuleName }

func (m *Module) WithAccountCache(cache AccountCache[string, int64]) *Module {
	if cache != nil {
		m.accCache = cache
	}

	return m
}
