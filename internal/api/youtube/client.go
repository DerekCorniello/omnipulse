// Package youtube provides a client for interacting with the YouTube Data API v3.
package youtube

import (
	"context"
	"net/http"

	"github.com/omnipulse/omnipulse/internal/config"
)

// Client handles communication with the YouTube Data API.
type Client struct {
	httpClient *http.Client
	config     config.YouTubeConfig
	baseURL    string
}

// NewClient creates a new YouTube API client.
func NewClient(cfg config.YouTubeConfig) *Client {
	return &Client{
		httpClient: &http.Client{},
		config:     cfg,
		baseURL:    "https://www.googleapis.com/youtube/v3",
	}
}

// SetHTTPClient allows setting a custom HTTP client (useful for testing).
func (c *Client) SetHTTPClient(client *http.Client) {
	c.httpClient = client
}

// GetChannelID returns the configured channel ID.
func (c *Client) GetChannelID() string {
	return c.config.ChannelID
}

// Close cleans up any resources held by the client.
func (c *Client) Close() error {
	// HTTP client doesn't need explicit cleanup
	return nil
}

// Ensure context is used (will be used in actual implementation)
var _ context.Context
