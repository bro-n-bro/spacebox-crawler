package bostrom

import (
	graph "github.com/cybercongress/go-cyber/x/graph/types"
)

func init() {
	registerTypes([]protoType{
		{(*graph.QueryGraphStatsRequest)(nil), "cyber.graph.v1beta1.QueryGraphStatsRequest"},
		{(*graph.QueryGraphStatsResponse)(nil), "cyber.graph.v1beta1.QueryGraphStatsResponse"},
		{(*graph.MsgCyberlink)(nil), "cyber.graph.v1beta1.MsgCyberlink"},
		{(*graph.MsgCyberlinkResponse)(nil), "cyber.graph.v1beta1.MsgCyberlinkResponse"},
		{(*graph.Link)(nil), "cyber.graph.v1beta1.Link"},
	})

	registerFiles([]protoFile{
		{"cyber/graph/v1beta1/tx.proto", descriptor((*graph.MsgCyberlink)(nil))},
		{"cyber/graph/v1beta1/query.proto", descriptor((*graph.QueryGraphStatsRequest)(nil))},
		{"cyber/graph/v1beta1/types.proto", descriptor((*graph.Link)(nil))},
	})
}
