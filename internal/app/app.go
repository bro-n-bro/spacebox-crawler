package app

import (
	"context"
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"

	ibctransfertypes "github.com/cosmos/ibc-go/v5/modules/apps/transfer/types"

	"github.com/cosmos/cosmos-sdk/simapp/params"

	"github.com/cosmos/cosmos-sdk/codec"
	cdc "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/std"
	cosmomodule "github.com/cosmos/cosmos-sdk/types/module"
	ibcstypes "github.com/cosmos/ibc-go/v5/modules/core/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	"bro-n-bro-osmosis/adapter/broker"
	grpcClient "bro-n-bro-osmosis/client/grpc"
	rpcClient "bro-n-bro-osmosis/client/rpc"
	"bro-n-bro-osmosis/internal/rep"
	"bro-n-bro-osmosis/modules"
	"bro-n-bro-osmosis/modules/messages"
	tb "bro-n-bro-osmosis/pkg/mapper/to_broker"
	"bro-n-bro-osmosis/pkg/worker"
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

func New(cfg Config) *App {
	l := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("cmp", "app").Logger()
	return &App{
		log: &l,
		cfg: cfg,
	}
}

func (a *App) Start(ctx context.Context) error {
	a.log.Info().Msg("starting app")

	grpcCli := grpcClient.New(a.cfg.GrpcConfig)
	rpcCli := rpcClient.New(a.cfg.RpcConfig)
	b := broker.New(a.cfg.BrokerConfig)

	//encoding := MakeEncodingConfig(getBasicManagers())
	cdc := MakeEncodingConfigV2()

	tb := tb.NewToBroker(cdc)
	parser := messages.JoinMessageParsers(messages.CosmosMessageAddressesParser)
	modules := modules.BuildModules(b, grpcCli, *tb, parser, cdc, a.cfg.Modules...)

	w := worker.New(a.cfg.WorkerConfig, b, rpcCli, grpcCli, modules, cdc, *tb)

	MakeSdkConfig(a.cfg, sdk.GetConfig())

	// TODO: mongo

	a.cmps = append(
		a.cmps,
		cmp{grpcCli, "grpcClient"},
		cmp{rpcCli, "rpcClient"},
		cmp{b, "broker"},
		cmp{w, "worker"},
	)

	okCh, errCh := make(chan struct{}), make(chan error)
	go func() {
		for _, c := range a.cmps {
			a.log.Info().Msgf("%v is starting", c.Name)
			if err := c.Service.Start(ctx); err != nil {
				a.log.Error().Err(err).Msgf(FmtCannotStart, c.Name)
				errCh <- errors.Wrapf(err, FmtCannotStart, c.Name)
			}
		}

		okCh <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ErrStartTimeout
	case err := <-errCh:
		return err
	case <-okCh:
		return nil
	}
}

func (a *App) Stop(ctx context.Context) error {
	a.log.Info().Msg("shutting down service...")

	errCh := make(chan error)
	go func() {
		gr, ctx := errgroup.WithContext(ctx)
		var c cmp
		for i := len(a.cmps) - 1; i >= 0; i-- {
			c = a.cmps[i]
			a.log.Info().Msgf("stopping %q...", c.Name)
			if err := c.Service.Stop(ctx); err != nil {
				a.log.Error().Err(err).Msgf("cannot stop %q", c.Name)
			}
		}
		errCh <- gr.Wait()
	}()

	select {
	case <-ctx.Done():
		return ErrShutdownTimeout
	case err := <-errCh:
		if err != nil {
			return err
		}
		return nil
	}
}

// TODO: move out it
// getBasicManagers returns the various basic managers that are used to register the encoding to
// support custom messages.
// This should be edited by custom implementations if needed.
func getBasicManagers() []cosmomodule.BasicManager {
	return []cosmomodule.BasicManager{
		simapp.ModuleBasics,
		//ibcsimapp.ModuleBasics,
		//{ibc.AppModuleBasic{}.Name(): ibc.AppModuleBasic{}},
		//{ibctransfer.AppModuleBasic{}.Name(): ibctransfer.AppModuleBasic{}},
		//{ibctmc.AppModuleBasic{}.Name(): ibctmc.AppModuleBasic{}},
	}
}

// MakeEncodingConfig creates an EncodingConfig to properly handle all the messages
func MakeEncodingConfig(managers []cosmomodule.BasicManager) params.EncodingConfig {
	encodingConfig := params.MakeTestEncodingConfig()

	//ibctypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	//ibctransfertypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	//ibctmctypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ibcstypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	manager := mergeBasicManagers(managers)
	manager.RegisterLegacyAminoCodec(encodingConfig.Amino)
	manager.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}

func MakeEncodingConfigV2() codec.Codec {

	ir := cdc.NewInterfaceRegistry()
	simapp.ModuleBasics.RegisterInterfaces(ir)
	std.RegisterInterfaces(ir)
	std.RegisterLegacyAminoCodec(codec.NewAminoCodec(codec.NewLegacyAmino()).LegacyAmino) // FIXME: not needed?
	ibcstypes.RegisterInterfaces(ir)
	ibctransfertypes.RegisterInterfaces(ir)

	return codec.NewProtoCodec(ir)
}

// mergeBasicManagers merges the given managers into a single module.BasicManager
func mergeBasicManagers(managers []cosmomodule.BasicManager) cosmomodule.BasicManager {
	var union = cosmomodule.BasicManager{}
	for _, manager := range managers {
		for k, v := range manager {
			union[k] = v
		}
	}
	return union
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
