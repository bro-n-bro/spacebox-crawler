package ibc

import (
	"os"
	"sync"

	"github.com/rs/zerolog"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	tb "github.com/bro-n-bro/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/bro-n-bro/spacebox-crawler/types"
)

const (
	moduleName = "ibc"
)

var (
	_ types.Module         = &Module{}
	_ types.MessageHandler = &Module{}
	_ types.BlockHandler   = &Module{}
)

type (
	denomCache struct {
		mu          sync.RWMutex
		denomHashes map[string]struct{}
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
	l := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("module", moduleName).Logger()

	return &Module{
		log:        &l,
		broker:     b,
		tbM:        tbM,
		client:     client,
		denomCache: &denomCache{denomHashes: make(map[string]struct{})},
	}
}

func (m *Module) Name() string { return moduleName }
