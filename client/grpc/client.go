package grpc

import (
	"context"
	"time"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/types/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"google.golang.org/grpc"
)

type Client struct {
	target string

	TmsService          tmservice.ServiceClient
	TxService           tx.ServiceClient
	BankQueryClient     banktypes.QueryClient
	GovQueryClient      govtypes.QueryClient
	MintQueryClient     minttypes.QueryClient
	SlashingQueryClient slashingtypes.QueryClient
	StakingQueryClient  stakingtypes.QueryClient

	conn *grpc.ClientConn
}

func New(target string) *Client {
	return &Client{target: target}
}

func (c *Client) Start(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// Create a connection to the gRPC server.
	grpcConn, err := grpc.DialContext(
		ctx,
		c.target,            // Or your gRPC server address.
		grpc.WithInsecure(), // The SDK doesn't support any transport security mechanism.
		grpc.WithBlock(),
	)
	if err != nil {
		return err
	}

	c.TmsService = tmservice.NewServiceClient(grpcConn)
	c.TxService = tx.NewServiceClient(grpcConn)
	c.BankQueryClient = banktypes.NewQueryClient(grpcConn)
	c.GovQueryClient = govtypes.NewQueryClient(grpcConn)
	c.MintQueryClient = minttypes.NewQueryClient(grpcConn)
	c.SlashingQueryClient = slashingtypes.NewQueryClient(grpcConn)
	c.StakingQueryClient = stakingtypes.NewQueryClient(grpcConn)

	c.conn = grpcConn
	return nil
}

func (c *Client) Stop(_ context.Context) error {
	return c.conn.Close()
}
