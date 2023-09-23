package app

import (
	"context"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	cdc "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/capability"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	groupmodule "github.com/cosmos/cosmos-sdk/x/group/module"
	"github.com/cosmos/cosmos-sdk/x/mint"
	nftmodule "github.com/cosmos/cosmos-sdk/x/nft/module"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibclightclient "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	interchainprovider "github.com/cosmos/interchain-security/v3/x/ccv/provider"
	interchaintypes "github.com/cosmos/interchain-security/v3/x/ccv/types"
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
	"github.com/bro-n-bro/spacebox-crawler/pkg/cache"
	tb "github.com/bro-n-bro/spacebox-crawler/pkg/mapper/to_broker"
	ts "github.com/bro-n-bro/spacebox-crawler/pkg/mapper/to_storage"
	"github.com/bro-n-bro/spacebox-crawler/pkg/worker"
	liquiditytypes "github.com/bro-n-bro/spacebox-crawler/types/liquidity"
)

const (
	defaultCacheSize = 100000
	FmtCannotStart   = "cannot start %q"
)

var (
	ErrStartTimeout    = errors.New("start timeout")
	ErrShutdownTimeout = errors.New("shutdown timeout")

	lessInt64    = func(cacheVal, newVal int64) bool { return cacheVal < newVal }
	greaterInt64 = func(cacheVal, newVal int64) bool { return cacheVal > newVal }
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

	// TODO: use redis
	valCache, err := cache.New[string, int64](defaultCacheSize, cache.WithCompareFunc[string, int64](lessInt64))
	if err != nil {
		return err
	}

	valCommissionCache, err := cache.New[string, int64](defaultCacheSize, cache.WithCompareFunc[string, int64](lessInt64))
	if err != nil {
		return err
	}

	valDescriptionCache, err := cache.New[string, int64](defaultCacheSize, cache.WithCompareFunc[string, int64](lessInt64)) //nolint:lll
	if err != nil {
		return err
	}

	valInfoCache, err := cache.New[string, int64](defaultCacheSize, cache.WithCompareFunc[string, int64](lessInt64))
	if err != nil {
		return err
	}

	valStatusCache, err := cache.New[string, int64](defaultCacheSize, cache.WithCompareFunc[string, int64](lessInt64))
	if err != nil {
		return err
	}

	// collect data only from bigger height
	tallyCache, err := cache.New[uint64, int64](defaultCacheSize, cache.WithCompareFunc[uint64, int64](lessInt64))
	if err != nil {
		return err
	}

	// collect data only from earlier height
	accCache, err := cache.New[string, int64](defaultCacheSize, cache.WithCompareFunc[string, int64](greaterInt64))
	if err != nil {
		return err
	}

	if a.cfg.MetricsEnabled {
		cache.RegisterMetrics("spacebox_crawler")

		valCache.Patch(cache.WithMetrics[string, int64]("validators"))
		valCommissionCache.Patch(cache.WithMetrics[string, int64]("validators_commission"))
		valDescriptionCache.Patch(cache.WithMetrics[string, int64]("validators_description"))
		valInfoCache.Patch(cache.WithMetrics[string, int64]("validators_info"))
		valStatusCache.Patch(cache.WithMetrics[string, int64]("validators_status"))
		tallyCache.Patch(cache.WithMetrics[uint64, int64]("tally"))
		accCache.Patch(cache.WithMetrics[string, int64]("account"))
	}

	b := broker.New(a.cfg.BrokerConfig, a.cfg.Modules, *a.log,
		broker.WithValidatorCache(valCache),
		broker.WithValidatorCommissionCache(valCommissionCache),
		broker.WithValidatorDescriptionCache(valDescriptionCache),
		broker.WithValidatorInfoCache(valInfoCache),
		broker.WithValidatorStatusCache(valStatusCache),
	)
	s := storage.New(a.cfg.StorageConfig, *a.log)

	cdc, amino := MakeEncodingConfig()
	tb := tb.NewToBroker(cdc, amino.LegacyAmino)
	parser := core.JoinMessageParsers(core.CosmosMessageAddressesParser)

	modules := modules.BuildModules(b, a.log, grpcCli, *tb, cdc, a.cfg.Modules, parser, a.cfg.DefaultDenom,
		tallyCache, accCache)

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

	var basicManager = module.NewBasicManager(
		auth.AppModuleBasic{},
		genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distribution.AppModuleBasic{},
		gov.NewAppModuleBasic(
			[]govclient.ProposalHandler{
				paramsclient.ProposalHandler,
				upgradeclient.LegacyProposalHandler,
				upgradeclient.LegacyCancelProposalHandler,
			},
		),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		feegrantmodule.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		authzmodule.AppModuleBasic{},
		groupmodule.AppModuleBasic{},
		vesting.AppModuleBasic{},
		nftmodule.AppModuleBasic{},
		consensus.AppModuleBasic{},
		ibc.AppModuleBasic{},
		ibclightclient.AppModuleBasic{},
		interchainprovider.AppModuleBasic{},
	)

	basicManager.RegisterInterfaces(ir)
	std.RegisterInterfaces(ir)
	ibctransfertypes.RegisterInterfaces(ir)
	cryptocodec.RegisterInterfaces(ir)
	interchaintypes.RegisterInterfaces(ir)
	liquiditytypes.RegisterInterfaces(ir)

	amino := codec.NewAminoCodec(codec.NewLegacyAmino())
	std.RegisterLegacyAminoCodec(amino.LegacyAmino) // FIXME: not needed?
	ibctransfertypes.RegisterLegacyAminoCodec(amino.LegacyAmino)
	liquiditytypes.RegisterLegacyAminoCodec(amino.LegacyAmino)

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
