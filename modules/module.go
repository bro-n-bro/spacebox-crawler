package modules

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	grpcClient "github.com/hexy-dev/spacebox-crawler/client/grpc"
	"github.com/hexy-dev/spacebox-crawler/internal/rep"
	authModule "github.com/hexy-dev/spacebox-crawler/modules/auth"
	bankModule "github.com/hexy-dev/spacebox-crawler/modules/bank"
	coreModule "github.com/hexy-dev/spacebox-crawler/modules/core"
	crisisModule "github.com/hexy-dev/spacebox-crawler/modules/crisis"
	distributionModule "github.com/hexy-dev/spacebox-crawler/modules/distribution"
	evidenceModule "github.com/hexy-dev/spacebox-crawler/modules/evidence"
	govModule "github.com/hexy-dev/spacebox-crawler/modules/gov"
	"github.com/hexy-dev/spacebox-crawler/modules/messages"
	mintModule "github.com/hexy-dev/spacebox-crawler/modules/mint"
	slashingModule "github.com/hexy-dev/spacebox-crawler/modules/slasing"
	stakingModule "github.com/hexy-dev/spacebox-crawler/modules/staking"
	tb "github.com/hexy-dev/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/hexy-dev/spacebox-crawler/types"
)

func BuildModules(b rep.Broker, cli *grpcClient.Client, tbMapper tb.ToBroker, addressesParser messages.MessageAddressesParser, cdc codec.Codec,
	modules []string) []types.Module {

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
		// module, ok := moduleStrMap[m]
		// if !ok {
		//	// todo: log
		//	continue
		// }
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
