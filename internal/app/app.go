package app

import (
	"context"
	"time"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	adminTypes "github.com/cosmos/admin-module/x/adminmodule/types"
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
	interchainprovider "github.com/cosmos/interchain-security/v4/x/ccv/provider"
	dmntypes "github.com/cybercongress/go-cyber/x/dmn/types"
	graphtypes "github.com/cybercongress/go-cyber/x/graph/types"
	gridtypes "github.com/cybercongress/go-cyber/x/grid/types"
	resourcestypes "github.com/cybercongress/go-cyber/x/resources/types"
	contractmanagertypes "github.com/neutron-org/neutron/v2/x/contractmanager/types"
	neutroncrontypes "github.com/neutron-org/neutron/v2/x/cron/types"
	neutrondextypes "github.com/neutron-org/neutron/v2/x/dex/types"
	neutronfeeburnertypes "github.com/neutron-org/neutron/v2/x/feeburner/types"
	neutronfeerefundertypes "github.com/neutron-org/neutron/v2/x/feerefunder/types"
	neutroninterchainqueriestypes "github.com/neutron-org/neutron/v2/x/interchainqueries/types"
	neutroninterchaintxstypes "github.com/neutron-org/neutron/v2/x/interchaintxs/types"
	neutrontokenfactorytypes "github.com/neutron-org/neutron/v2/x/tokenfactory/types"
	neutrontransfertypes "github.com/neutron-org/neutron/v2/x/transfer/types"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog"
	blocksdktypes "github.com/skip-mev/block-sdk/x/auction/types"

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
	_ "github.com/bro-n-bro/spacebox-crawler/types/bostrom"
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
		log     *zerolog.Logger
		version string
		cmps    []cmp
		cfg     Config
	}
	cmp struct {
		Service rep.Lifecycle
		Name    string
	}
)

func New(cfg Config, version string, l zerolog.Logger) *App {
	l = l.With().Str("version", version).Str("cmp", "app").Logger()

	return &App{
		log:     &l,
		cfg:     cfg,
		version: version,
	}
}

func (a *App) Start(ctx context.Context) error {
	a.log.Info().Msg("starting app")

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
	tCache, err := cache.New[uint64, int64](defaultCacheSize, cache.WithCompareFunc[uint64, int64](lessInt64))
	if err != nil {
		return err
	}

	// collect data only from earlier height
	aCache, err := cache.New[string, int64](defaultCacheSize, cache.WithCompareFunc[string, int64](greaterInt64))
	if err != nil {
		return err
	}

	// collect data only from bigger height
	rCache, err := cache.New[string, int64](defaultCacheSize, cache.WithCompareFunc[string, int64](lessInt64))
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
		tCache.Patch(cache.WithMetrics[uint64, int64]("tally"))
		aCache.Patch(cache.WithMetrics[string, int64]("account"))

		promauto.NewGauge(prometheus.GaugeOpts{
			Namespace:   "spacebox_crawler",
			Name:        "version",
			Help:        "Crawler version",
			ConstLabels: prometheus.Labels{"version": a.version},
		}).Inc()
	}

	var (
		cod, amn = MakeEncodingConfig()
		sto      = storage.New(a.cfg.StorageConfig, *a.log)
		rpcCli   = rpcClient.New(a.cfg.RPCConfig)
		grpcCli  = grpcClient.New(a.cfg.GRPCConfig, *a.log, sto)
		tbr      = tb.NewToBroker(cod, amn.LegacyAmino)
		par      = core.JoinMessageParsers(core.CosmosMessageAddressesParser)

		brk = broker.New(a.cfg.BrokerConfig, a.cfg.Modules, *a.log,
			broker.WithValidatorCache(valCache),
			broker.WithValidatorCommissionCache(valCommissionCache),
			broker.WithValidatorDescriptionCache(valDescriptionCache),
			broker.WithValidatorInfoCache(valInfoCache),
			broker.WithValidatorStatusCache(valStatusCache),
		)

		mds = modules.BuildModules(a.log, a.cfg.Modules, a.cfg.DefaultDenom, grpcCli, rpcCli, brk, cod, *tbr, par, tCache, aCache, rCache) //nolint:lll
		tos = ts.NewToStorage()
		wrk = worker.New(a.cfg.WorkerConfig, *a.log, brk, rpcCli, grpcCli, mds, sto, cod, *tbr, *tos)
		srv = server.New(a.cfg.Server, sto, *a.log)
	)

	MakeSDKConfig(a.cfg, sdk.GetConfig())

	a.cmps = append(a.cmps,
		cmp{sto, "storage"},
		cmp{grpcCli, "grpc_client"},
		cmp{rpcCli, "rpc_client"},
		cmp{brk, "broker"},
		cmp{wrk, "worker"},
		cmp{srv, "server"},
	)

	okCh, errCh := make(chan struct{}), make(chan error)

	go func() {
		for _, c := range a.cmps {
			a.log.Info().Str("service", c.Name).Msg("starting")

			if err := c.Service.Start(ctx); err != nil {
				a.log.Error().Err(err).Msgf(FmtCannotStart, c.Name)
				errCh <- errors.Wrapf(err, FmtCannotStart, c.Name)

				return
			}

			a.log.Info().Str("service", c.Name).Msg("started")
		}

		okCh <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ErrStartTimeout
	case err := <-errCh:
		return err
	case <-okCh:
		a.log.Info().Msg("application started")
		return nil
	}
}

func (a *App) Stop(ctx context.Context) error {
	a.log.Info().Msg("shutting down service...")

	okCh, errCh := make(chan struct{}), make(chan error)

	go func() {
		for i := len(a.cmps) - 1; i > 0; i-- {
			c := a.cmps[i]

			a.log.Info().Str("service", c.Name).Msg("stopping")

			if err := c.Service.Stop(ctx); err != nil {
				a.log.Error().Str("service", c.Name).Err(err).Msg("cannot stop")
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
		a.log.Info().Msg("application stopped")
		return nil
	}
}

func (a *App) GetStartTimeout() time.Duration { return a.cfg.StartTimeout }
func (a *App) GetStopTimeout() time.Duration  { return a.cfg.StopTimeout }

// MakeEncodingConfig creates an EncodingConfig to properly handle and marshal all messages
func MakeEncodingConfig() (codec.Codec, *codec.AminoCodec) {
	var (
		registry     = cdc.NewInterfaceRegistry()
		basicManager = module.NewBasicManager(
			auth.AppModuleBasic{},
			genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
			bank.AppModuleBasic{},
			capability.AppModuleBasic{},
			staking.AppModuleBasic{},
			mint.AppModuleBasic{},
			distribution.AppModuleBasic{},
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
			gov.NewAppModuleBasic(
				[]govclient.ProposalHandler{
					paramsclient.ProposalHandler,
					upgradeclient.LegacyProposalHandler,
					upgradeclient.LegacyCancelProposalHandler,
				},
			),
		)
	)

	//
	basicManager.RegisterInterfaces(registry)
	std.RegisterInterfaces(registry)
	ibctransfertypes.RegisterInterfaces(registry)
	cryptocodec.RegisterInterfaces(registry)
	liquiditytypes.RegisterInterfaces(registry)

	// bostrom
	graphtypes.RegisterInterfaces(registry)
	dmntypes.RegisterInterfaces(registry)
	gridtypes.RegisterInterfaces(registry)
	resourcestypes.RegisterInterfaces(registry)
	wasmtypes.RegisterInterfaces(registry)

	// neutron
	adminTypes.RegisterInterfaces(registry)
	contractmanagertypes.RegisterInterfaces(registry)
	neutroncrontypes.RegisterInterfaces(registry)
	neutrondextypes.RegisterInterfaces(registry)
	neutronfeeburnertypes.RegisterInterfaces(registry)
	neutronfeerefundertypes.RegisterInterfaces(registry)
	neutroninterchainqueriestypes.RegisterInterfaces(registry)
	neutroninterchaintxstypes.RegisterInterfaces(registry)
	neutrontokenfactorytypes.RegisterInterfaces(registry)
	neutrontransfertypes.RegisterInterfaces(registry)
	blocksdktypes.RegisterInterfaces(registry)

	//
	amino := codec.NewAminoCodec(codec.NewLegacyAmino())
	std.RegisterLegacyAminoCodec(amino.LegacyAmino) // FIXME: not needed?
	ibctransfertypes.RegisterLegacyAminoCodec(amino.LegacyAmino)
	liquiditytypes.RegisterLegacyAminoCodec(amino.LegacyAmino)

	return codec.NewProtoCodec(registry), amino
}

// MakeSDKConfig represents a handy implementation of SdkConfigSetup that simply setups the prefix
// inside the configuration
func MakeSDKConfig(cfg Config, sdkConfig *sdk.Config) {
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
