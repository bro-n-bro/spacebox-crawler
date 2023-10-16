package grpc

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	feegranttypes "github.com/cosmos/cosmos-sdk/x/feegrant"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	bandwidthtypes "github.com/cybercongress/go-cyber/x/bandwidth/types"
	dmntypes "github.com/cybercongress/go-cyber/x/dmn/types"
	graphtypes "github.com/cybercongress/go-cyber/x/graph/types"
	gridtypes "github.com/cybercongress/go-cyber/x/grid/types"
	ranktypes "github.com/cybercongress/go-cyber/x/rank/types"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	liquiditytypes "github.com/bro-n-bro/spacebox-crawler/types/liquidity"
)

type (
	Client struct {
		SlashingQueryClient     slashingtypes.QueryClient
		TmsService              tmservice.ServiceClient
		TxService               tx.ServiceClient
		BankQueryClient         banktypes.QueryClient
		AuthQueryClient         authtypes.QueryClient
		GovQueryClient          govtypes.QueryClient
		MintQueryClient         minttypes.QueryClient
		StakingQueryClient      stakingtypes.QueryClient
		DistributionQueryClient distributiontypes.QueryClient
		AuthzQueryClient        authztypes.QueryClient
		FeegrantQueryClient     feegranttypes.QueryClient
		IbcTransferQueryClient  ibctransfertypes.QueryClient
		LiquidityQueryClient    liquiditytypes.QueryClient
		GraphQueryClient        graphtypes.QueryClient
		BandwidthQueryClient    bandwidthtypes.QueryClient
		DMNQueryClient          dmntypes.QueryClient
		GridQueryClient         gridtypes.QueryClient
		RankQueryClient         ranktypes.QueryClient
		conn                    *grpc.ClientConn
		cfg                     Config
	}
)

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
	c.AuthzQueryClient = authztypes.NewQueryClient(grpcConn)
	c.FeegrantQueryClient = feegranttypes.NewQueryClient(grpcConn)
	c.IbcTransferQueryClient = ibctransfertypes.NewQueryClient(grpcConn)
	c.LiquidityQueryClient = liquiditytypes.NewQueryClient(grpcConn)
	c.AuthQueryClient = authtypes.NewQueryClient(grpcConn)
	c.GraphQueryClient = graphtypes.NewQueryClient(grpcConn)
	c.BandwidthQueryClient = bandwidthtypes.NewQueryClient(grpcConn)
	c.DMNQueryClient = dmntypes.NewQueryClient(grpcConn)
	c.GridQueryClient = gridtypes.NewQueryClient(grpcConn)
	c.RankQueryClient = ranktypes.NewQueryClient(grpcConn)

	c.conn = grpcConn

	return nil
}

func (c *Client) Stop(_ context.Context) error { return c.conn.Close() }
