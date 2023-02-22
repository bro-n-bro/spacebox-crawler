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

	RPCClient  *tmHttp.HTTP
	HTTPClient *http.Client

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

	c.HTTPClient = httpCli

	// FIXME: does not work without websocket connection
	var rpcCli *tmHttp.HTTP
	if c.cfg.WSEnabled {
		rpcCli, err = tmHttp.NewWithClient(c.cfg.Host, "/websocket", httpCli)
		if err != nil {
			return err
		}

		if err = rpcCli.Start(); err != nil {
			return err
		}
	} else {
		rpcCli, err = tmHttp.NewWithClient(c.cfg.Host, "", httpCli)
		if err != nil {
			return err
		}
		if err = rpcCli.Start(); err != nil {
			return err
		}
	}

	c.RPCClient = rpcCli

	return nil
}

func (c *Client) Stop(ctx context.Context) error {
	c.HTTPClient.CloseIdleConnections()

	if c.cfg.WSEnabled {
		if err := c.RPCClient.Stop(); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) WsEnabled() bool { return c.cfg.WSEnabled }
