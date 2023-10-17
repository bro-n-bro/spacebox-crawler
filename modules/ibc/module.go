package ibc

import (
	"sync"

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
	_ types.Module         = &Module{}
	_ types.MessageHandler = &Module{}
	_ types.BlockHandler   = &Module{}
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
		denomCache *denomCache
	}
)

func New(b broker, tbM tb.ToBroker, client *grpcClient.Client) *Module {
	return &Module{
		log:        utils.NewModuleLogger(ModuleName),
		broker:     b,
		tbM:        tbM,
		client:     client,
		denomCache: &denomCache{denomHashes: make(map[string]struct{})},
	}
}

func (m *Module) Name() string { return ModuleName }
