package app

import (
	"context"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	cdc "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v5/modules/apps/transfer/types"
	ibcstypes "github.com/cosmos/ibc-go/v5/modules/core/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/bro-n-bro/spacebox-crawler/adapter/storage"
	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	rpcClient "github.com/bro-n-bro/spacebox-crawler/client/rpc"
	"github.com/bro-n-bro/spacebox-crawler/delivery/broker"
	"github.com/bro-n-bro/spacebox-crawler/delivery/server"
	"github.com/bro-n-bro/spacebox-crawler/internal/rep"
	"github.com/bro-n-bro/spacebox-crawler/modules"
	"github.com/bro-n-bro/spacebox-crawler/modules/core"
	tb "github.com/bro-n-bro/spacebox-crawler/pkg/mapper/to_broker"
	ts "github.com/bro-n-bro/spacebox-crawler/pkg/mapper/to_storage"
	"github.com/bro-n-bro/spacebox-crawler/pkg/worker"
)

const (
	FmtCannotStart = "cannot start %q"
)

var (
	ErrStartTimeout    = errors.New("start timeout")
	ErrShutdownTimeout = errors.New("shutdown timeout")
)

type (
	App struct {
		log  *zerolog.Logger
		cmps []cmp
		cfg  Config
	}
	cmp struct {
		Service rep.Lifecycle
		Name    string
	}
)

func New(cfg Config, l zerolog.Logger) *App {
	l = l.With().Str("cmp", "app").Logger()

	return &App{
		log: &l,
		cfg: cfg,
	}
}

func (a *App) Start(ctx context.Context) error {
	a.log.Info().Msg("starting app")

	grpcCli := grpcClient.New(a.cfg.GRPCConfig)
	rpcCli := rpcClient.New(a.cfg.RPCConfig)
	b := broker.New(a.cfg.BrokerConfig, a.cfg.Modules, *a.log)
	s := storage.New(a.cfg.StorageConfig, *a.log)

	// encoding := MakeEncodingConfig(getBasicManagers())
	cdc, amino := MakeEncodingConfig()

	tb := tb.NewToBroker(cdc, amino.LegacyAmino)
	parser := core.JoinMessageParsers(core.CosmosMessageAddressesParser)
	modules := modules.BuildModules(b, grpcCli, *tb, cdc, a.cfg.Modules, parser)
	ts := ts.NewToStorage()

	w := worker.New(a.cfg.WorkerConfig, *a.log, b, rpcCli, grpcCli, modules, s, cdc, *tb, *ts)
	server := server.New(a.cfg.Server, s, *a.log)

	MakeSdkConfig(a.cfg, sdk.GetConfig())

	a.cmps = append(
		a.cmps,
		cmp{s, "storage"},
		cmp{grpcCli, "grpcClient"},
		cmp{rpcCli, "rpcClient"},
		cmp{b, "broker"},
		cmp{w, "worker"},
		cmp{server, "server"},
	)

	okCh, errCh := make(chan struct{}), make(chan error)

	go func() {
		for _, c := range a.cmps {
			a.log.Info().Msgf("%v is starting", c.Name)

			if err := c.Service.Start(ctx); err != nil {
				a.log.Error().Err(err).Msgf(FmtCannotStart, c.Name)
				errCh <- errors.Wrapf(err, FmtCannotStart, c.Name)

				return
			}

			a.log.Info().Msgf("%v started", c.Name)
		}
		okCh <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ErrStartTimeout
	case err := <-errCh:
		return err
	case <-okCh:
		a.log.Info().Msg("Application started!")
		return nil
	}
}

func (a *App) Stop(ctx context.Context) error {
	a.log.Info().Msg("shutting down service...")

	okCh, errCh := make(chan struct{}), make(chan error)

	go func() {
		for i := len(a.cmps) - 1; i > 0; i-- {
			c := a.cmps[i]
			a.log.Info().Msgf("stopping %q...", c.Name)

			if err := c.Service.Stop(ctx); err != nil {
				a.log.Error().Err(err).Msgf("cannot stop %q", c.Name)
				errCh <- err

				return
			}
		}

		okCh <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ErrShutdownTimeout
	case err := <-errCh:
		return err
	case <-okCh:
		a.log.Info().Msg("Application stopped!")
		return nil
	}
}

func (a *App) GetStartTimeout() time.Duration { return a.cfg.StartTimeout }
func (a *App) GetStopTimeout() time.Duration  { return a.cfg.StopTimeout }

// MakeEncodingConfig creates an EncodingConfig to properly handle and marshal all messages
func MakeEncodingConfig() (codec.Codec, *codec.AminoCodec) {
	ir := cdc.NewInterfaceRegistry()

	simapp.ModuleBasics.RegisterInterfaces(ir)
	std.RegisterInterfaces(ir)
	ibcstypes.RegisterInterfaces(ir)
	ibctransfertypes.RegisterInterfaces(ir)

	amino := codec.NewAminoCodec(codec.NewLegacyAmino())
	std.RegisterLegacyAminoCodec(amino.LegacyAmino) // FIXME: not needed?
	ibctransfertypes.RegisterLegacyAminoCodec(amino.LegacyAmino)
	// liquiditytypes.RegisterInterfaces(ir)

	return codec.NewProtoCodec(ir), amino
}

// MakeSdkConfig represents a handy implementation of SdkConfigSetup that simply setups the prefix
// inside the configuration
func MakeSdkConfig(cfg Config, sdkConfig *sdk.Config) {
	prefix := cfg.ChainPrefix
	sdkConfig.SetBech32PrefixForAccount(
		prefix,
		prefix+sdk.PrefixPublic,
	)
	sdkConfig.SetBech32PrefixForValidator(
		prefix+sdk.PrefixValidator+sdk.PrefixOperator,
		prefix+sdk.PrefixValidator+sdk.PrefixOperator+sdk.PrefixPublic,
	)
	sdkConfig.SetBech32PrefixForConsensusNode(
		prefix+sdk.PrefixValidator+sdk.PrefixConsensus,
		prefix+sdk.PrefixValidator+sdk.PrefixConsensus+sdk.PrefixPublic,
	)
}
