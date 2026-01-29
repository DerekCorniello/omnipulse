// Package linkedin provides a client for interacting with the LinkedIn API.
package linkedin

import (
	"context"
	"net/http"

	"github.com/crossforge/crossforge/internal/config"
)

// Client handles communication with the LinkedIn API.
type Client struct {
	httpClient *http.Client
	config     config.LinkedInConfig
	baseURL    string
}

// NewClient creates a new LinkedIn API client.
func NewClient(cfg config.LinkedInConfig) *Client {
	return &Client{
		httpClient: &http.Client{},
		config:     cfg,
		baseURL:    "https://api.linkedin.com/v2",
	}
}

// SetHTTPClient allows setting a custom HTTP client (useful for testing).
func (c *Client) SetHTTPClient(client *http.Client) {
	c.httpClient = client
}

// GetPersonURN returns the configured person URN.
func (c *Client) GetPersonURN() string {
	return c.config.PersonURN
}

// Close cleans up any resources held by the client.
func (c *Client) Close() error {
	return nil
}

// Ensure context is used
var _ context.Context
