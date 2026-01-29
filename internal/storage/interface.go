// Package storage provides data persistence interfaces and implementations.
package storage

import (
	"context"

	"github.com/omnipulse/omnipulse/internal/data"
)

// Store defines the interface for data persistence operations.
type Store interface {
	// Video operations
	SaveVideo(ctx context.Context, video *data.Video) error
	GetVideo(ctx context.Context, id string) (*data.Video, error)
	GetVideos(ctx context.Context, limit, offset int) ([]*data.Video, error)
	GetVideosByDateRange(ctx context.Context, dateRange data.DateRange) ([]*data.Video, error)

	// Tweet operations
	SaveTweet(ctx context.Context, tweet *data.Tweet) error
	GetTweet(ctx context.Context, id string) (*data.Tweet, error)
	GetTweets(ctx context.Context, limit, offset int) ([]*data.Tweet, error)
	GetTweetsByDateRange(ctx context.Context, dateRange data.DateRange) ([]*data.Tweet, error)

	// LinkedIn post operations
	SaveLinkedInPost(ctx context.Context, post *data.LinkedInPost) error
	GetLinkedInPost(ctx context.Context, id string) (*data.LinkedInPost, error)
	GetLinkedInPosts(ctx context.Context, limit, offset int) ([]*data.LinkedInPost, error)
	GetLinkedInPostsByDateRange(ctx context.Context, dateRange data.DateRange) ([]*data.LinkedInPost, error)

	// Comment operations
	SaveComment(ctx context.Context, comment *data.Comment) error
	GetComments(ctx context.Context, platform data.Platform, contentID string) ([]*data.Comment, error)

	// Channel/User stats operations
	SaveChannelStats(ctx context.Context, stats *data.ChannelStats) error
	GetLatestChannelStats(ctx context.Context) (*data.ChannelStats, error)
	SaveXUserStats(ctx context.Context, stats *data.XUserStats) error
	GetLatestXUserStats(ctx context.Context) (*data.XUserStats, error)
	SaveLinkedInProfileStats(ctx context.Context, stats *data.LinkedInProfileStats) error
	GetLatestLinkedInProfileStats(ctx context.Context) (*data.LinkedInProfileStats, error)

	// Insight operations
	SaveInsight(ctx context.Context, insight *data.Insight) error
	GetInsights(ctx context.Context, platform data.Platform, limit int) ([]*data.Insight, error)
	GetRecentInsights(ctx context.Context, limit int) ([]*data.Insight, error)

	// Analytics summary operations
	GetAnalyticsSummary(ctx context.Context, dateRange data.DateRange) (*data.AnalyticsSummary, error)

	// Trend operations
	GetTrendData(ctx context.Context, platform data.Platform, metric string, days int) (*data.TrendData, error)

	// Database management
	Migrate(ctx context.Context) error
	Close() error
}
