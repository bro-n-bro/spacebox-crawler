package modules

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"bro-n-bro-osmosis/adapter/broker"
	grpcClient "bro-n-bro-osmosis/client/grpc"
	authModule "bro-n-bro-osmosis/modules/auth"
	bankModule "bro-n-bro-osmosis/modules/bank"
	govModule "bro-n-bro-osmosis/modules/gov"
	"bro-n-bro-osmosis/modules/messages"
	mintModule "bro-n-bro-osmosis/modules/mint"
	slashingModule "bro-n-bro-osmosis/modules/slasing"
	stakingModule "bro-n-bro-osmosis/modules/staking"
	"bro-n-bro-osmosis/types"
)

var moduleStrMap = map[string]types.Module{
	"bank": &bankModule.Module{},
}

func BuildModules(b *broker.Broker, cli *grpcClient.Client, parser messages.MessageAddressesParser, cdc codec.Codec,
	modules ...string) []types.Module {

	res := make([]types.Module, 0)
	for _, m := range modules {
		// TODO: make better
		switch m {
		case "auth":
			res = append(res, authModule.New(b, cli, cdc, parser))
		case "bank":
			res = append(res, bankModule.New(b, cli, cdc, parser))
		case "gov":
			res = append(res, govModule.New(b, cli, cdc))
		case "mint":
			res = append(res, mintModule.New(b, cli))
		case "slashing":
			res = append(res, slashingModule.New(b, cli))
		case "staking":
			res = append(res, stakingModule.New(b, cli, cdc, parser, modules))

		default:
			continue
		}
		//module, ok := moduleStrMap[m]
		//if !ok {
		//	// todo: log
		//	continue
		//}
	}
	return res
}
