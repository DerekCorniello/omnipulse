// Package x provides a client for interacting with the X (Twitter) API v2.
package x

import (
	"context"
	"net/http"

	"github.com/crossforge/crossforge/internal/config"
)

// Client handles communication with the X (Twitter) API v2.
type Client struct {
	httpClient  *http.Client
	config      config.XConfig
	baseURL     string
}

// NewClient creates a new X API client.
func NewClient(cfg config.XConfig) *Client {
	return &Client{
		httpClient: &http.Client{},
		config:     cfg,
		baseURL:    "https://api.twitter.com/2",
	}
}

// SetHTTPClient allows setting a custom HTTP client (useful for testing).
func (c *Client) SetHTTPClient(client *http.Client) {
	c.httpClient = client
}

// GetUserID returns the configured user ID.
func (c *Client) GetUserID() string {
	return c.config.UserID
}

// GetBearerToken returns the configured bearer token.
func (c *Client) GetBearerToken() string {
	return c.config.BearerToken
}

// Close cleans up any resources held by the client.
func (c *Client) Close() error {
	return nil
}

// Ensure context is used
var _ context.Context
