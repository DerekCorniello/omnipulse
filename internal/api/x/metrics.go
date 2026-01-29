// Package x provides metrics fetching for X (Twitter) accounts.
package x

import (
	"context"

	"github.com/crossforge/crossforge/internal/data"
)

// Metrics handles fetching X metrics and analytics.
type Metrics struct {
	client *Client
}

// NewMetrics creates a new Metrics instance.
func NewMetrics(client *Client) *Metrics {
	return &Metrics{client: client}
}

// GetUserStats fetches user-level statistics.
// Uses X API v2: GET /2/users/:id with user.fields
func (m *Metrics) GetUserStats(ctx context.Context) (*data.XUserStats, error) {
	// TODO: Implement using X API v2
	// Endpoint: GET https://api.twitter.com/2/users/{id}
	// Headers:
	//   - Authorization: Bearer {bearer_token}
	// Parameters:
	//   - user.fields: public_metrics,created_at,description,profile_image_url
	//
	// Response includes:
	//   - public_metrics.followers_count
	//   - public_metrics.following_count
	//   - public_metrics.tweet_count
	//   - public_metrics.listed_count
	return nil, nil
}

// GetUserTweets fetches recent tweets from the user.
// Uses X API v2: GET /2/users/:id/tweets
func (m *Metrics) GetUserTweets(ctx context.Context, maxResults int) ([]*data.Tweet, error) {
	// TODO: Implement using X API v2
	// Endpoint: GET https://api.twitter.com/2/users/{id}/tweets
	// Headers:
	//   - Authorization: Bearer {bearer_token}
	// Parameters:
	//   - max_results: 10-100
	//   - tweet.fields: public_metrics,created_at,conversation_id,in_reply_to_user_id
	//   - expansions: attachments.media_keys (optional)
	//   - exclude: retweets,replies (optional)
	//
	// Response includes:
	//   - id, text, created_at
	//   - public_metrics.retweet_count
	//   - public_metrics.reply_count
	//   - public_metrics.like_count
	//   - public_metrics.quote_count
	//   - public_metrics.impression_count (requires OAuth 2.0 User Context)
	return nil, nil
}

// GetTweetMetrics fetches detailed metrics for a specific tweet.
// Uses X API v2: GET /2/tweets/:id
func (m *Metrics) GetTweetMetrics(ctx context.Context, tweetID string) (*data.Tweet, error) {
	// TODO: Implement using X API v2
	// Endpoint: GET https://api.twitter.com/2/tweets/{id}
	// Headers:
	//   - Authorization: Bearer {bearer_token}
	// Parameters:
	//   - tweet.fields: public_metrics,non_public_metrics,organic_metrics,created_at
	//
	// Note: non_public_metrics and organic_metrics require OAuth 2.0 User Context
	// and are only available for tweets authored by the authenticated user
	//
	// Non-public metrics include:
	//   - impression_count
	//   - url_link_clicks
	//   - user_profile_clicks
	return nil, nil
}

// GetMentions fetches recent mentions of the user.
// Uses X API v2: GET /2/users/:id/mentions
func (m *Metrics) GetMentions(ctx context.Context, maxResults int) ([]*data.Tweet, error) {
	// TODO: Implement using X API v2
	// Endpoint: GET https://api.twitter.com/2/users/{id}/mentions
	// Headers:
	//   - Authorization: Bearer {bearer_token}
	// Parameters:
	//   - max_results: 10-100
	//   - tweet.fields: public_metrics,created_at,author_id
	//   - expansions: author_id
	//   - user.fields: name,username,profile_image_url
	return nil, nil
}
