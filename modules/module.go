package modules

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/rs/zerolog"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	"github.com/bro-n-bro/spacebox-crawler/client/rpc"
	"github.com/bro-n-bro/spacebox-crawler/internal/rep"
	authModule "github.com/bro-n-bro/spacebox-crawler/modules/auth"
	authzModule "github.com/bro-n-bro/spacebox-crawler/modules/authz"
	bandwidthModule "github.com/bro-n-bro/spacebox-crawler/modules/bandwidth"
	bankModule "github.com/bro-n-bro/spacebox-crawler/modules/bank"
	coreModule "github.com/bro-n-bro/spacebox-crawler/modules/core"
	distributionModule "github.com/bro-n-bro/spacebox-crawler/modules/distribution"
	dmnModule "github.com/bro-n-bro/spacebox-crawler/modules/dmn"
	feeGrantModule "github.com/bro-n-bro/spacebox-crawler/modules/feegrant"
	govModule "github.com/bro-n-bro/spacebox-crawler/modules/gov"
	graphModule "github.com/bro-n-bro/spacebox-crawler/modules/graph"
	gridModule "github.com/bro-n-bro/spacebox-crawler/modules/grid"
	ibcModule "github.com/bro-n-bro/spacebox-crawler/modules/ibc"
	liquidityModule "github.com/bro-n-bro/spacebox-crawler/modules/liquidity"
	mintModule "github.com/bro-n-bro/spacebox-crawler/modules/mint"
	rankModule "github.com/bro-n-bro/spacebox-crawler/modules/rank"
	rawModule "github.com/bro-n-bro/spacebox-crawler/modules/raw"
	resourcesModule "github.com/bro-n-bro/spacebox-crawler/modules/resources"
	slashingModule "github.com/bro-n-bro/spacebox-crawler/modules/slashing"
	stakingModule "github.com/bro-n-bro/spacebox-crawler/modules/staking"
	wasmModule "github.com/bro-n-bro/spacebox-crawler/modules/wasm"
	tb "github.com/bro-n-bro/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/bro-n-bro/spacebox-crawler/types"
)

type (
	Cache[K, V comparable] interface{ UpdateCacheValue(K, V) bool }

	tallyCache   Cache[uint64, int64]
	accountCache Cache[string, int64]
	routeCache   Cache[string, int64]
)

//nolint:gocyclo
func BuildModules(
	log *zerolog.Logger,
	mds []string,
	den string,
	cli *grpcClient.Client,
	rpc *rpc.Client,
	brk rep.Broker,
	cdc codec.Codec,
	tbm tb.ToBroker,
	aParse coreModule.MsgAddrParser,
	tCache tallyCache,
	aCache accountCache,
	rCache routeCache,
) []types.Module {

	mods := NewModuleLoader().WithLogger(log)

	for _, mod := range mds {
		switch mod {
		case authModule.ModuleName:
			auth := authModule.New(brk, cli, tbm, cdc, aParse).WithAccountCache(aCache)
			mods.Add(auth)
		case bankModule.ModuleName:
			bank := bankModule.New(brk, cli, tbm, cdc, aParse)
			mods.Add(bank)
		case govModule.ModuleName:
			gov := govModule.New(brk, cli, tbm, cdc).WithTallyCache(tCache)
			mods.Add(gov)
		case mintModule.ModuleName:
			mint := mintModule.New(brk, cli, tbm)
			mods.Add(mint)
		case stakingModule.ModuleName:
			staking := stakingModule.New(brk, cli, tbm, cdc, mds, den).WithAccountCache(aCache)
			mods.Add(staking)
		case distributionModule.ModuleName:
			distribution := distributionModule.New(brk, cli, rpc, tbm, cdc)
			mods.Add(distribution)
		case coreModule.ModuleName:
			core := coreModule.New(brk, tbm, cdc, aParse)
			mods.Add(core)
		case authzModule.ModuleName:
			authz := authzModule.New(brk, cli, tbm, cdc)
			mods.Add(authz)
		case feeGrantModule.ModuleName:
			feeGrant := feeGrantModule.New(brk, cli, tbm, cdc)
			mods.Add(feeGrant)
		case slashingModule.ModuleName:
			slashing := slashingModule.New(brk, cli, tbm)
			mods.Add(slashing)
		case ibcModule.ModuleName:
			ibc := ibcModule.New(brk, tbm, cli, cdc)
			mods.Add(ibc)
		case liquidityModule.ModuleName:
			liquidity := liquidityModule.New(brk, cli, tbm)
			mods.Add(liquidity)
		case graphModule.ModuleName:
			graph := graphModule.New(brk, tbm, cdc, cli)
			mods.Add(graph)
		case bandwidthModule.ModuleName:
			bandwidth := bandwidthModule.New(brk, cli, tbm)
			mods.Add(bandwidth)
		case dmnModule.ModuleName:
			dmn := dmnModule.New(brk, cli, tbm)
			mods.Add(dmn)
		case gridModule.ModuleName:
			grid := gridModule.New(brk, cli, tbm).WithCache(rCache)
			mods.Add(grid)
		case rankModule.ModuleName:
			rank := rankModule.New(brk, cli, tbm)
			mods.Add(rank)
		case resourcesModule.ModuleName:
			resources := resourcesModule.New(brk, cli, tbm)
			mods.Add(resources)
		case wasmModule.ModuleName:
			wasm := wasmModule.New(brk, cdc)
			mods.Add(wasm)
		case rawModule.ModuleName:
			raw := rawModule.New(brk, rpc, tbm)
			mods.Add(raw)
		default:
			log.Warn().Str("name", mod).Msg("unknown module")
			continue
		}
	}

	return mods.Build()
}
