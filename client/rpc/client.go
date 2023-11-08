package rpc

import (
	"context"

	cometbftHttp "github.com/cometbft/cometbft/rpc/client/http"
	jsonrpcclient "github.com/cometbft/cometbft/rpc/jsonrpc/client"
)

type Client struct {
	*jsonrpcclient.WSClient
	*cometbftHttp.WSEvents

	RPCClient *cometbftHttp.HTTP

	cfg Config
}

func New(cfg Config) *Client {
	return &Client{cfg: cfg}
}

func (c *Client) Start(ctx context.Context) error {
	// FIXME: does not work without websocket connection
	var rpcCli *cometbftHttp.HTTP
	if c.cfg.WSEnabled {
		var err error
		rpcCli, err = cometbftHttp.NewWithTimeout(c.cfg.Host, "/websocket", 15)
		if err != nil {
			return err
		}

		if err = rpcCli.Start(); err != nil {
			return err
		}
	} else {
		var err error
		rpcCli, err = cometbftHttp.NewWithTimeout(c.cfg.Host, "", 15)
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
	if c.cfg.WSEnabled {
		if err := c.RPCClient.Stop(); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) WsEnabled() bool { return c.cfg.WSEnabled }
