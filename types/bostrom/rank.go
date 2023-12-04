package bostrom

import (
	ranktypes "github.com/cybercongress/go-cyber/x/rank/types"
)

func init() {
	registerTypes([]protoType{
		{(*ranktypes.Params)(nil), "cyber.rank.v1beta1.Params"},
		{(*ranktypes.RankedParticle)(nil), "cyber.rank.v1beta1.RankedParticle"},

		{(*ranktypes.QueryParamsRequest)(nil), "cyber.rank.v1beta1.QueryParamsRequest"},
		{(*ranktypes.QueryParamsResponse)(nil), "cyber.rank.v1beta1.QueryParamsResponse"},
		{(*ranktypes.QueryRankRequest)(nil), "cyber.rank.v1beta1.QueryRankRequest"},
		{(*ranktypes.QueryRankResponse)(nil), "cyber.rank.v1beta1.QueryRankResponse"},
		{(*ranktypes.QuerySearchRequest)(nil), "cyber.rank.v1beta1.QuerySearchRequest"},
		{(*ranktypes.QuerySearchResponse)(nil), "cyber.rank.v1beta1.QuerySearchResponse"},
		{(*ranktypes.QueryTopRequest)(nil), "cyber.rank.v1beta1.QueryTopRequest"},
		{(*ranktypes.QueryIsLinkExistRequest)(nil), "cyber.rank.v1beta1.QueryIsLinkExistRequest"},
		{(*ranktypes.QueryIsAnyLinkExistRequest)(nil), "cyber.rank.v1beta1.QueryIsAnyLinkExistRequest"},
		{(*ranktypes.QueryLinkExistResponse)(nil), "cyber.rank.v1beta1.QueryLinkExistResponse"},
		{(*ranktypes.QueryNegentropyPartilceRequest)(nil), "cyber.rank.v1beta1.QueryNegentropyPartilceRequest"},
		{(*ranktypes.QueryNegentropyParticleResponse)(nil), "cyber.rank.v1beta1.QueryNegentropyParticleResponse"},
		{(*ranktypes.QueryNegentropyRequest)(nil), "cyber.rank.v1beta1.QueryNegentropyRequest"},
		{(*ranktypes.QueryNegentropyResponse)(nil), "cyber.rank.v1beta1.QueryNegentropyResponse"},
		{(*ranktypes.QueryKarmaRequest)(nil), "cyber.rank.v1beta1.QueryKarmaRequest"},
		{(*ranktypes.QueryKarmaResponse)(nil), "cyber.rank.v1beta1.QueryKarmaResponse"},

		{(*ranktypes.GenesisState)(nil), "cyber.rank.v1beta1.GenesisState"},
	})

	registerFiles([]protoFile{
		{"cyber/rank/v1beta1/types.proto", descriptor((*ranktypes.Params)(nil))},
		{"cyber/rank/v1beta1/query.proto", descriptor((*ranktypes.QueryParamsRequest)(nil))},
		{"cyber/rank/v1beta1/genesis.proto", descriptor((*ranktypes.GenesisState)(nil))},
	})
}
