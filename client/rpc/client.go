package rpc

import (
	"context"
	"net/http"

	tmHttp "github.com/tendermint/tendermint/rpc/client/http"
	jsonrpcclient "github.com/tendermint/tendermint/rpc/jsonrpc/client"
)

type Client struct {
	*jsonrpcclient.WSClient
	*tmHttp.WSEvents

	RpcClient  *tmHttp.HTTP
	HttpClient *http.Client

	cfg Config
}

func New(cfg Config) *Client {
	return &Client{cfg: cfg}
}

func (c *Client) Start(ctx context.Context) error {
	httpCli, err := jsonrpcclient.DefaultHTTPClient(c.cfg.Host)
	if err != nil {
		return err
	}

	c.HttpClient = httpCli

	// FIXME: does not work without websocket
	if c.cfg.WSEnabled {
		cli, err := tmHttp.NewWithClient(c.cfg.Host, "/websocket", httpCli)
		if err != nil {
			return err
		}
		if err = cli.Start(); err != nil {
			return err
		}
		c.RpcClient = cli
	} else {
		//return nil
		cli, err := tmHttp.NewWithClient(c.cfg.Host, "", httpCli)
		if err != nil {
			return err
		}
		if err = cli.Start(); err != nil {
			return err
		}
		c.RpcClient = cli
	}

	return nil
}

func (c *Client) Stop(ctx context.Context) error {
	c.HttpClient.CloseIdleConnections()
	if c.cfg.WSEnabled {
		if err := c.RpcClient.Stop(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) WsEnabled() bool { return c.cfg.WSEnabled }
