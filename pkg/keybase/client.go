package keybase

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

var (
	// DefaultHTTPClient is the default rate limited HTTP client
	DefaultHTTPClient = &RLHTTPClient{
		client:      http.DefaultClient,
		Ratelimiter: rate.NewLimiter(rate.Every(1*time.Second), 25), // 25 request per 1 second
	}
)

// RLHTTPClient Rate Limited HTTP Client
type RLHTTPClient struct {
	client      *http.Client
	Ratelimiter *rate.Limiter
}

// Do dispatches the HTTP request to the network
func (c *RLHTTPClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	// Comment out the below 3 lines to turn off rate limiting
	if err := c.Ratelimiter.Wait(ctx); err != nil { // This is a blocking call. Honours the rate limit
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
