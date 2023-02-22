package grpc

import (
	"strconv"

	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// GetHeightRequestHeader returns the grpc.CallOption to query the state at a given height
func GetHeightRequestHeader(height int64) grpc.CallOption {
	header := metadata.New(map[string]string{
		grpctypes.GRPCBlockHeightHeader: strconv.FormatInt(height, 10),
	})

	return grpc.Header(&header)
}
