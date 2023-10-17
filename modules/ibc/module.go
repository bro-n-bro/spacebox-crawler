package ibc

import (
	"sync"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/rs/zerolog"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	"github.com/bro-n-bro/spacebox-crawler/modules/utils"
	tb "github.com/bro-n-bro/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/bro-n-bro/spacebox-crawler/types"
)

const (
	ModuleName = "ibc"
)

var (
	_ types.Module                   = &Module{}
	_ types.RecursiveMessagesHandler = &Module{}
	_ types.BlockHandler             = &Module{}
)

type (
	denomCache struct {
		denomHashes map[string]struct{}
		mu          sync.RWMutex
	}

	Module struct {
		log        *zerolog.Logger
		client     *grpcClient.Client
		broker     broker
		tbM        tb.ToBroker
		cdc        codec.Codec
		denomCache *denomCache
	}
)

func New(b broker, tbM tb.ToBroker, client *grpcClient.Client, cdc codec.Codec) *Module {
	return &Module{
		log:        utils.NewModuleLogger(ModuleName),
		broker:     b,
		tbM:        tbM,
		client:     client,
		cdc:        cdc,
		denomCache: &denomCache{denomHashes: make(map[string]struct{})},
	}
}

func (m *Module) Name() string { return ModuleName }
