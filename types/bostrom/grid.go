package bostrom

import (
	gridtypes "github.com/cybercongress/go-cyber/x/grid/types"
)

func init() {

	registerTypes([]protoType{
		{(*gridtypes.Params)(nil), "cyber.grid.v1beta1.Params"},
		{(*gridtypes.Route)(nil), "cyber.grid.v1beta1.Route"},
		{(*gridtypes.Value)(nil), "cyber.grid.v1beta1.Value"},

		{(*gridtypes.MsgCreateRoute)(nil), "cyber.grid.v1beta1.MsgCreateRoute"},
		{(*gridtypes.MsgEditRoute)(nil), "cyber.grid.v1beta1.MsgEditRoute"},
		{(*gridtypes.MsgDeleteRoute)(nil), "cyber.grid.v1beta1.MsgDeleteRoute"},
		{(*gridtypes.MsgEditRouteName)(nil), "cyber.grid.v1beta1.MsgEditRouteName"},
		{(*gridtypes.MsgCreateRouteResponse)(nil), "cyber.grid.v1beta1.MsgCreateRouteResponse"},
		{(*gridtypes.MsgEditRouteResponse)(nil), "cyber.grid.v1beta1.MsgEditRouteResponse"},
		{(*gridtypes.MsgDeleteRouteResponse)(nil), "cyber.grid.v1beta1.MsgDeleteRouteResponse"},
		{(*gridtypes.MsgEditRouteNameResponse)(nil), "cyber.grid.v1beta1.MsgEditRouteNameResponse"},

		{(*gridtypes.QueryParamsRequest)(nil), "cyber.grid.v1beta1.QueryParamsRequest"},
		{(*gridtypes.QueryParamsResponse)(nil), "cyber.grid.v1beta1.QueryParamsResponse"},
		{(*gridtypes.QuerySourceRequest)(nil), "cyber.grid.v1beta1.QuerySourceRequest"},
		{(*gridtypes.QueryDestinationRequest)(nil), "cyber.grid.v1beta1.QueryDestinationRequest"},
		{(*gridtypes.QueryRoutedEnergyResponse)(nil), "cyber.grid.v1beta1.QueryRoutedEnergyResponse"},
		{(*gridtypes.QueryRouteRequest)(nil), "cyber.grid.v1beta1.QueryRouteRequest"},
		{(*gridtypes.QueryRouteResponse)(nil), "cyber.grid.v1beta1.QueryRouteResponse"},
		{(*gridtypes.QueryRoutesRequest)(nil), "cyber.grid.v1beta1.QueryRoutesRequest"},
		{(*gridtypes.QueryRoutesResponse)(nil), "cyber.grid.v1beta1.QueryRoutesResponse"},

		{(*gridtypes.GenesisState)(nil), "cyber.grid.v1beta1.GenesisState"},
	})

	registerFiles([]protoFile{
		{"cyber/grid/v1beta1/types.proto", descriptor((*gridtypes.Params)(nil))},
		{"cyber/grid/v1beta1/tx.proto", descriptor((*gridtypes.MsgCreateRoute)(nil))},
		{"cyber/grid/v1beta1/query.proto", descriptor((*gridtypes.QueryParamsRequest)(nil))},
		{"cyber/grid/v1beta1/genesis.proto", descriptor((*gridtypes.GenesisState)(nil))},
	})
}
