package rpc

import (
	"context"
	"net/http"

	tmHttp "github.com/tendermint/tendermint/rpc/client/http"
	"github.com/tendermint/tendermint/rpc/jsonrpc/client"
)

type Client struct {
	target    string
	wsEnabled bool

	RpcClient  *tmHttp.HTTP
	HttpClient *http.Client
}

func New(target string, wsEnabled bool) *Client {
	return &Client{target: target, wsEnabled: wsEnabled}
}

func (c *Client) Start(ctx context.Context) error {
	httpCli, err := client.DefaultHTTPClient(c.target)
	if err != nil {
		return err
	}

	if c.wsEnabled {
		cli, err := tmHttp.NewWithClient(c.target, "/websocket", httpCli)
		if err != nil {
			return err
		}
		if err = cli.Start(); err != nil {
			return err
		}
		c.RpcClient = cli
	}

	c.HttpClient = httpCli

	return nil
}

func (c *Client) Stop(ctx context.Context) error {
	c.HttpClient.CloseIdleConnections()
	if c.wsEnabled {
		if err := c.RpcClient.Stop(); err != nil {
			return err
		}
	}
	return nil
}
