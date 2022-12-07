package modules

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	grpcClient "bro-n-bro-osmosis/client/grpc"
	"bro-n-bro-osmosis/internal/rep"
	authModule "bro-n-bro-osmosis/modules/auth"
	bankModule "bro-n-bro-osmosis/modules/bank"
	coreModule "bro-n-bro-osmosis/modules/core"
	crisisModule "bro-n-bro-osmosis/modules/crisis"
	distributionModule "bro-n-bro-osmosis/modules/distribution"
	evidenceModule "bro-n-bro-osmosis/modules/evidence"
	govModule "bro-n-bro-osmosis/modules/gov"
	"bro-n-bro-osmosis/modules/messages"
	mintModule "bro-n-bro-osmosis/modules/mint"
	slashingModule "bro-n-bro-osmosis/modules/slasing"
	stakingModule "bro-n-bro-osmosis/modules/staking"
	tb "bro-n-bro-osmosis/pkg/mapper/to_broker"
	"bro-n-bro-osmosis/types"
)

var moduleStrMap = map[string]types.Module{
	"bank": &bankModule.Module{},
}

func BuildModules(b rep.Broker, cli *grpcClient.Client, tbMapper tb.ToBroker, addressesParser messages.MessageAddressesParser, cdc codec.Codec,
	modules ...string) []types.Module {

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
		case "slashing":
			res = append(res, slashingModule.New(b, cli))
		case "staking":
			res = append(res, stakingModule.New(b, cli, tbMapper, cdc, modules))
		case "distribution":
			res = append(res, distributionModule.New(b, cli, tbMapper, cdc, addressesParser))
		case "core":
			res = append(res, coreModule.New(b, tbMapper, cdc, addressesParser))

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

func BuildAddressesParser(modules ...string) AddressesParser {
	parsers := []AddressesParser{DefaultMessagesParser}
	for _, m := range modules {
		// TODO: make better
		switch m {
		case "bank":
			parsers = append(parsers, bankModule.BankAccountsParser)
		case "gov":
			parsers = append(parsers, govModule.GovMessagesParser)
		case "slashing":
			parsers = append(parsers, slashingModule.SlashingMessagesParser)
		case "staking":
			parsers = append(parsers, stakingModule.StakingMessagesParser)
		case "evidence":
			parsers = append(parsers, evidenceModule.EvidenceMessagesParser)
		case "crisis":
			parsers = append(parsers, crisisModule.CrisisMessagesParser)
		case "distribution":

		default:
			continue
		}
		//module, ok := moduleStrMap[m]
		//if !ok {
		//	// todo: log
		//	continue
		//}
	}

	return func(cdc codec.Codec, msg sdk.Msg) ([]string, error) {
		for _, parser := range parsers {
			if accounts, _ := parser(cdc, msg); len(accounts) > 0 {
				return accounts, nil
			}

		}
		return nil, nil
	}
}
