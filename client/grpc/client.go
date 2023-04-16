package grpc

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/types/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	feegranttypes "github.com/cosmos/cosmos-sdk/x/feegrant"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	SlashingQueryClient     slashingtypes.QueryClient
	TmsService              tmservice.ServiceClient
	TxService               tx.ServiceClient
	BankQueryClient         banktypes.QueryClient
	GovQueryClient          govtypes.QueryClient
	MintQueryClient         minttypes.QueryClient
	StakingQueryClient      stakingtypes.QueryClient
	DistributionQueryClient distributiontypes.QueryClient
	FeegrantQueryClient     feegranttypes.QueryClient
	conn                    *grpc.ClientConn
	cfg                     Config
}

func New(cfg Config) *Client {
	return &Client{cfg: cfg}
}

func (c *Client) Start(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	options := []grpc.DialOption{
		grpc.WithBlock(),
	}

	if c.cfg.MetricsEnabled {
		options = append(
			options,
			grpc.WithUnaryInterceptor(grpc_prometheus.UnaryClientInterceptor),
			grpc.WithStreamInterceptor(grpc_prometheus.StreamClientInterceptor),
		)

		grpc_prometheus.EnableClientHandlingTimeHistogram()
	}

	// Add required secure grpc option based on config parameter
	if c.cfg.SecureConnection {
		options = append(options, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{}))) // nolint:gosec
	} else {
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Create a connection to the gRPC server.
	grpcConn, err := grpc.DialContext(
		ctx,
		c.cfg.Host, // Or your gRPC server address.
		options...,
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
	c.DistributionQueryClient = distributiontypes.NewQueryClient(grpcConn)

	c.conn = grpcConn

	return nil
}

func (c *Client) Stop(_ context.Context) error {
	return c.conn.Close()
}
