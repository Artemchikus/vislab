package gitlab

import (
	"net/http"
	"time"
)

type ClientOption func(*Client) error

func WithAPIPrefix(apiPrefix string) ClientOption {
	return func(c *Client) error {
		return c.setApiPrefix(apiPrefix)
	}
}

func WithRateLimit(ms int64) ClientOption {
	return func(c *Client) error {
		c.rateLimit = time.Duration(ms) * time.Millisecond
		return nil
	}
}

func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) error {
		c.client = client
		return nil
	}
}

func WithPerPage(perPage int64) ClientOption {
	return func(c *Client) error {
		c.perPage = perPage
		return nil
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) error {
		c.client.Timeout = timeout
		return nil
	}
}
