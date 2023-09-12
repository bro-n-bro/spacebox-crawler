package ibc

import (
	"os"
	"sync"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/rs/zerolog"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	tb "github.com/bro-n-bro/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/bro-n-bro/spacebox-crawler/types"
)

const (
	moduleName = "ibc"
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
	l := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("module", moduleName).Logger()

	return &Module{
		log:        &l,
		broker:     b,
		tbM:        tbM,
		client:     client,
		cdc:        cdc,
		denomCache: &denomCache{denomHashes: make(map[string]struct{})},
	}
}

func (m *Module) Name() string { return moduleName }
