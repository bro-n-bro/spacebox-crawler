package bostrom

import (
	dmntypes "github.com/cybercongress/go-cyber/x/dmn/types"
)

func init() {
	registerTypes([]protoType{
		{(*dmntypes.MsgCreateThought)(nil), "cyber.dmn.v1beta1.MsgCreateThought"},
		{(*dmntypes.MsgForgetThought)(nil), "cyber.dmn.v1beta1.MsgForgetThought"},
		{(*dmntypes.MsgChangeThoughtParticle)(nil), "cyber.dmn.v1beta1.MsgChangeThoughtParticle"},
		{(*dmntypes.MsgChangeThoughtInput)(nil), "cyber.dmn.v1beta1.MsgChangeThoughtInput"},
		{(*dmntypes.MsgChangeThoughtName)(nil), "cyber.dmn.v1beta1.MsgChangeThoughtName"},
		{(*dmntypes.MsgChangeThoughtGasPrice)(nil), "cyber.dmn.v1beta1.MsgChangeThoughtGasPrice"},
		{(*dmntypes.MsgChangeThoughtPeriod)(nil), "cyber.dmn.v1beta1.MsgChangeThoughtPeriod"},
		{(*dmntypes.MsgChangeThoughtBlock)(nil), "cyber.dmn.v1beta1.MsgChangeThoughtBlock"},
		{(*dmntypes.MsgCreateThoughtResponse)(nil), "cyber.dmn.v1beta1.MsgCreateThoughtResponse"},
		{(*dmntypes.MsgForgetThoughtResponse)(nil), "cyber.dmn.v1beta1.MsgForgetThoughtResponse"},
		{(*dmntypes.MsgChangeThoughtParticleResponse)(nil), "cyber.dmn.v1beta1.MsgChangeThoughtParticleResponse"},
		{(*dmntypes.MsgChangeThoughtNameResponse)(nil), "cyber.dmn.v1beta1.MsgChangeThoughtNameResponse"},
		{(*dmntypes.MsgChangeThoughtInputResponse)(nil), "cyber.dmn.v1beta1.MsgChangeThoughtInputResponse"},
		{(*dmntypes.MsgChangeThoughtGasPriceResponse)(nil), "cyber.dmn.v1beta1.MsgChangeThoughtGasPriceResponse"},
		{(*dmntypes.MsgChangeThoughtPeriodResponse)(nil), "cyber.dmn.v1beta1.MsgChangeThoughtPeriodResponse"},
		{(*dmntypes.MsgChangeThoughtBlockResponse)(nil), "cyber.dmn.v1beta1.MsgChangeThoughtBlockResponse"},

		{(*dmntypes.Params)(nil), "cyber.dmn.v1beta1.Params"},
		{(*dmntypes.Thought)(nil), "cyber.dmn.v1beta1.Thought"},
		{(*dmntypes.Trigger)(nil), "cyber.dmn.v1beta1.Trigger"},
		{(*dmntypes.Load)(nil), "cyber.dmn.v1beta1.Load"},
		{(*dmntypes.ThoughtStats)(nil), "cyber.dmn.v1beta1.ThoughtStats"},

		{(*dmntypes.QueryParamsRequest)(nil), "cyber.dmn.v1beta1.QueryParamsRequest"},
		{(*dmntypes.QueryParamsResponse)(nil), "cyber.dmn.v1beta1.QueryParamsResponse"},
		{(*dmntypes.QueryThoughtParamsRequest)(nil), "cyber.dmn.v1beta1.QueryThoughtParamsRequest"},
		{(*dmntypes.QueryThoughtResponse)(nil), "cyber.dmn.v1beta1.QueryThoughtResponse"},
		{(*dmntypes.QueryThoughtStatsResponse)(nil), "cyber.dmn.v1beta1.QueryThoughtStatsResponse"},
		{(*dmntypes.QueryThoughtsRequest)(nil), "cyber.dmn.v1beta1.QueryThoughtsRequest"},
		{(*dmntypes.QueryThoughtsResponse)(nil), "cyber.dmn.v1beta1.QueryThoughtsResponse"},
		{(*dmntypes.QueryThoughtsStatsRequest)(nil), "cyber.dmn.v1beta1.QueryThoughtsStatsRequest"},
		{(*dmntypes.QueryThoughtsStatsResponse)(nil), "cyber.dmn.v1beta1.QueryThoughtsStatsResponse"},

		{(*dmntypes.GenesisState)(nil), "cyber.dmn.v1beta1.GenesisState"},
	})

	registerFiles([]protoFile{
		{"cyber/dmn/v1beta1/tx.proto", descriptor((*dmntypes.MsgCreateThought)(nil))},
		{"cyber/dmn/v1beta1/types.proto", descriptor((*dmntypes.Params)(nil))},
		{"cyber/dmn/v1beta1/query.proto", descriptor((*dmntypes.QueryParamsRequest)(nil))},
		{"cyber/dmn/v1beta1/genesis.proto", descriptor((*dmntypes.GenesisState)(nil))},
	})
}
