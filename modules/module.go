package modules

import (
	"github.com/cosmos/cosmos-sdk/codec"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	"github.com/bro-n-bro/spacebox-crawler/internal/rep"
	authModule "github.com/bro-n-bro/spacebox-crawler/modules/auth"
	bankModule "github.com/bro-n-bro/spacebox-crawler/modules/bank"
	coreModule "github.com/bro-n-bro/spacebox-crawler/modules/core"
	distributionModule "github.com/bro-n-bro/spacebox-crawler/modules/distribution"
	govModule "github.com/bro-n-bro/spacebox-crawler/modules/gov"
	mintModule "github.com/bro-n-bro/spacebox-crawler/modules/mint"
	stakingModule "github.com/bro-n-bro/spacebox-crawler/modules/staking"
	tb "github.com/bro-n-bro/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/bro-n-bro/spacebox-crawler/types"
)

func BuildModules(b rep.Broker, cli *grpcClient.Client, tbMapper tb.ToBroker, cdc codec.Codec, modules []string,
	addressesParser coreModule.MessageAddressesParser) []types.Module {

	res := make([]types.Module, 0)

	for _, m := range modules {
		// TODO: make better
		switch m {
		case "auth":
			res = append(res, authModule.New(b, cli, tbMapper, cdc, addressesParser))
		case "bank":
			res = append(res, bankModule.New(b, cli, tbMapper, cdc, addressesParser))
		case "gov":
			res = append(res, govModule.New(b, cli, tbMapper, cdc))
		case "mint":
			res = append(res, mintModule.New(b, cli, tbMapper))
		case "staking":
			s := stakingModule.New(b, cli, tbMapper, cdc, modules)
			res = append(res, s)
		case "distribution":
			res = append(res, distributionModule.New(b, cli, tbMapper, cdc, addressesParser))
		case "core":
			res = append(res, coreModule.New(b, tbMapper, cdc, addressesParser))
		default:
			// TODO: log
			continue
		}
	}

	return res
}
