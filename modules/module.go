package modules

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/rs/zerolog"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	"github.com/bro-n-bro/spacebox-crawler/internal/rep"
	authModule "github.com/bro-n-bro/spacebox-crawler/modules/auth"
	authzModule "github.com/bro-n-bro/spacebox-crawler/modules/authz"
	bankModule "github.com/bro-n-bro/spacebox-crawler/modules/bank"
	coreModule "github.com/bro-n-bro/spacebox-crawler/modules/core"
	distributionModule "github.com/bro-n-bro/spacebox-crawler/modules/distribution"
	feegrantModule "github.com/bro-n-bro/spacebox-crawler/modules/feegrant"
	govModule "github.com/bro-n-bro/spacebox-crawler/modules/gov"
	ibcModule "github.com/bro-n-bro/spacebox-crawler/modules/ibc"
	mintModule "github.com/bro-n-bro/spacebox-crawler/modules/mint"
	slashingModule "github.com/bro-n-bro/spacebox-crawler/modules/slashing"
	stakingModule "github.com/bro-n-bro/spacebox-crawler/modules/staking"
	tb "github.com/bro-n-bro/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/bro-n-bro/spacebox-crawler/types"
)

func BuildModules(b rep.Broker, log *zerolog.Logger, cli *grpcClient.Client, tbMapper tb.ToBroker,
	cdc codec.Codec, modules []string, addressesParser coreModule.MessageAddressesParser, parseAvatarURL bool,
	tallyCache govModule.TallyCache[uint64, int64]) []types.Module {

	res := make([]types.Module, 0)

	for _, m := range modules {
		// TODO: make better
		switch m {
		case "auth":
			log.Info().Msg("auth module registered")
			res = append(res, authModule.New(b, cli, tbMapper, cdc, addressesParser))
		case "bank":
			log.Info().Msg("bank module registered")
			res = append(res, bankModule.New(b, cli, tbMapper, cdc, addressesParser))
		case "gov":
			gov := govModule.New(b, cli, tbMapper, cdc)
			if tallyCache != nil {
				gov.SetTallyCache(tallyCache)
			}
			log.Info().Msg("gov module registered")
			res = append(res, gov)
		case "mint":
			log.Info().Msg("mint module registered")
			res = append(res, mintModule.New(b, cli, tbMapper))
		case "staking":
			log.Info().Msg("staking module registered")
			res = append(res, stakingModule.New(b, cli, tbMapper, cdc, modules, parseAvatarURL))
		case "distribution":
			log.Info().Msg("distribution module registered")
			res = append(res, distributionModule.New(b, cli, tbMapper, cdc))
		case "core":
			log.Info().Msg("core module registered")
			res = append(res, coreModule.New(b, tbMapper, cdc, addressesParser))
		case "authz":
			log.Info().Msg("authz module registered")
			res = append(res, authzModule.New(b, cli, tbMapper, cdc))
		case "feegrant":
			log.Info().Msg("feegrant module registered")
			res = append(res, feegrantModule.New(b, cli, tbMapper, cdc))
		case "slashing":
			log.Info().Msg("slashing module registered")
			res = append(res, slashingModule.New(b, tbMapper))
		case "ibc":
			log.Info().Msg("ibc module registered")
			res = append(res, ibcModule.New(b, tbMapper))
		default:
			// TODO: log
			log.Warn().Msgf("unknown module: %v", m)
			continue
		}
	}

	return res
}
