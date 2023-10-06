package bostrom

import (
	resourcestypes "github.com/cybercongress/go-cyber/x/resources/types"
)

func init() {

	registerTypes([]protoType{
		{(*resourcestypes.Params)(nil), "cyber.resources.v1beta1.Params"},

		{(*resourcestypes.MsgInvestmint)(nil), "cyber.resources.v1beta1.MsgInvestmint"},
		{(*resourcestypes.MsgInvestmintResponse)(nil), "cyber.resources.v1beta1.MsgInvestmintResponse"},

		{(*resourcestypes.QueryParamsRequest)(nil), "cyber.resources.v1beta1.QueryParamsRequest"},
		{(*resourcestypes.QueryParamsResponse)(nil), "cyber.resources.v1beta1.QueryParamsResponse"},
		{(*resourcestypes.QueryInvestmintRequest)(nil), "cyber.resources.v1beta1.QueryInvestmintRequest"},
		{(*resourcestypes.QueryInvestmintResponse)(nil), "cyber.resources.v1beta1.QueryInvestmintResponse"},

		{(*resourcestypes.GenesisState)(nil), "cyber.resources.v1beta1.GenesisState"},
	})

	registerFiles([]protoFile{
		{"cyber/resources/v1beta1/types.proto", descriptor((*resourcestypes.Params)(nil))},
		{"cyber/resources/v1beta1/tx.proto", descriptor((*resourcestypes.MsgInvestmint)(nil))},
		{"cyber/resources/v1beta1/query.proto", descriptor((*resourcestypes.QueryParamsRequest)(nil))},
		{"cyber/resources/v1beta1/genesis.proto", descriptor((*resourcestypes.GenesisState)(nil))},
	})
}
