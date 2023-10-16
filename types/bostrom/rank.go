package bostrom

import (
	ranktypes "github.com/cybercongress/go-cyber/x/rank/types"
)

func init() { // TODO: @malekviktor i need a help here
	registerTypes([]protoType{
		{(*ranktypes.Params)(nil), "cyber.rank.v1beta1.Params"},

		{(*ranktypes.QueryParamsRequest)(nil), "cyber.rank.v1beta1.QueryParamsRequest"},
		{(*ranktypes.QueryParamsResponse)(nil), "cyber.rank.v1beta1.QueryParamsResponse"},

		{(*ranktypes.GenesisState)(nil), "cyber.rank.v1beta1.GenesisState"},
	})

	registerFiles([]protoFile{
		{"cyber/rank/v1beta1/types.proto", descriptor((*ranktypes.Params)(nil))},
		{"cyber/rank/v1beta1/query.proto", descriptor((*ranktypes.QueryParamsRequest)(nil))},
		{"cyber/rank/v1beta1/genesis.proto", descriptor((*ranktypes.GenesisState)(nil))},
	})
}
