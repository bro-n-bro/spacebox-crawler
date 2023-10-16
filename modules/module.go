package modules

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/rs/zerolog"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
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
	slashingModule "github.com/bro-n-bro/spacebox-crawler/modules/slashing"
	stakingModule "github.com/bro-n-bro/spacebox-crawler/modules/staking"
	tb "github.com/bro-n-bro/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/bro-n-bro/spacebox-crawler/types"
)

type (
	Cache[K, V comparable] interface{ UpdateCacheValue(K, V) bool }

	tallyCache   Cache[uint64, int64]
	accountCache Cache[string, int64]
)

func BuildModules(
	brk rep.Broker,
	log *zerolog.Logger,
	cli *grpcClient.Client,
	tbm tb.ToBroker,
	cdc codec.Codec,
	mds []string,
	aParse coreModule.MsgAddrParser,
	tCache tallyCache,
	aCache accountCache,
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
			staking := stakingModule.New(brk, cli, tbm, cdc, mds).WithAccountCache(aCache)
			mods.Add(staking)
		case distributionModule.ModuleName:
			distribution := distributionModule.New(brk, cli, tbm, cdc)
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
			ibc := ibcModule.New(brk, tbm, cli)
			mods.Add(ibc)
		case liquidityModule.ModuleName:
			liquidity := liquidityModule.New(brk, cli, tbm)
			mods.Add(liquidity)
		case graphModule.ModuleName:
			graph := graphModule.New(brk, tbm, cdc, cli)
			mods.Add(graph) // TODO: add to env vars
		case bandwidthModule.ModuleName:
			bandwidth := bandwidthModule.New(brk, cli, tbm)
			mods.Add(bandwidth) // TODO: add to env vars
		case dmnModule.ModuleName:
			dmn := dmnModule.New(brk, cli, tbm)
			mods.Add(dmn) // TODO: add to env vars
		case gridModule.ModuleName:
			grid := gridModule.New(brk, cli, tbm)
			mods.Add(grid) // TODO: add to env vars
		case rankModule.ModuleName:
			rank := rankModule.New(brk, cli, tbm)
			mods.Add(rank) // TODO: add to env vars
		default:
			log.Warn().Msgf("unknown module: %v", mod)
			continue
		}
	}

	return mods.Build()
}
