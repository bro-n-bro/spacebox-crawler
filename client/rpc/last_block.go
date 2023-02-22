package rpc

import "context"

func (c *Client) GetLastBlockHeight(ctx context.Context) (int64, error) {
	resp, err := c.RPCClient.ABCIInfo(ctx)
	if err != nil {
		return 0, err
	}

	return resp.Response.LastBlockHeight, nil
}
