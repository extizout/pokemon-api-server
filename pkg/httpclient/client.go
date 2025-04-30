package httpclient

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/labstack/gommon/log"
)

type Client struct {
	httpClient *http.Client
	retryCfg   RetryConfig
}

func NewClient(timeout time.Duration, retryCfg RetryConfig) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		retryCfg: retryCfg,
	}
}

func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	response, err := DoWithRetry(ctx, c.httpClient, req, c.retryCfg)
	if err != nil {
		switch {
		case errors.Is(err, ErrRequestTimeout):
			log.Warn("Timed out while making request")
		case errors.Is(err, ErrConnectionError):
			log.Warn("Network connection failed")
		default:
			log.Error("Unexpected HTTP error:", err)
		}
		return nil, err
	}
	return response, nil
}

func (c *Client) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(ctx, req)
}

func (c *Client) Post(ctx context.Context, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	return c.Do(ctx, req)
}
