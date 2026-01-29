// Package insights provides analytics aggregation and insight generation.
package insights

import (
	"context"

	"github.com/crossforge/crossforge/internal/data"
	"github.com/crossforge/crossforge/internal/storage"
)

// Aggregator aggregates analytics data across platforms.
type Aggregator struct {
	store storage.Store
}

// NewAggregator creates a new Aggregator.
func NewAggregator(store storage.Store) *Aggregator {
	return &Aggregator{store: store}
}

// GetDashboardData retrieves aggregated data for the dashboard view.
func (a *Aggregator) GetDashboardData(ctx context.Context, days int) (*DashboardData, error) {
	// TODO: Implement aggregation logic
	// 1. Get analytics summary for the date range
	// 2. Get recent insights
	// 3. Calculate trends for key metrics
	// 4. Get top performing content from each platform
	return nil, nil
}

// DashboardData contains all data needed for the main dashboard.
type DashboardData struct {
	Summary       *data.AnalyticsSummary `json:"summary"`
	RecentInsights []*data.Insight       `json:"recent_insights"`
	YouTubeTrend  *data.TrendData        `json:"youtube_trend,omitempty"`
	XTrend        *data.TrendData        `json:"x_trend,omitempty"`
	LinkedInTrend *data.TrendData        `json:"linkedin_trend,omitempty"`
	TopContent    *TopContent            `json:"top_content"`
}

// TopContent holds top performing content from each platform.
type TopContent struct {
	YouTube  []*data.Video        `json:"youtube,omitempty"`
	X        []*data.Tweet        `json:"x,omitempty"`
	LinkedIn []*data.LinkedInPost `json:"linkedin,omitempty"`
}

// GetPlatformAnalytics retrieves detailed analytics for a specific platform.
func (a *Aggregator) GetPlatformAnalytics(ctx context.Context, platform data.Platform, days int) (*PlatformAnalytics, error) {
	// TODO: Implement platform-specific aggregation
	return nil, nil
}

// PlatformAnalytics contains detailed analytics for a single platform.
type PlatformAnalytics struct {
	Platform    data.Platform   `json:"platform"`
	Summary     interface{}     `json:"summary"` // Platform-specific summary
	Trends      []*data.TrendData `json:"trends"`
	Content     interface{}     `json:"content"` // Platform-specific content list
	Insights    []*data.Insight `json:"insights"`
}

// ComparePlatforms compares performance metrics across platforms.
func (a *Aggregator) ComparePlatforms(ctx context.Context, days int) (*PlatformComparison, error) {
	// TODO: Implement cross-platform comparison
	// Normalize metrics across platforms (e.g., engagement rate)
	// Identify which platform is performing best
	return nil, nil
}

// PlatformComparison holds comparative analytics across platforms.
type PlatformComparison struct {
	BestPerforming     data.Platform            `json:"best_performing"`
	EngagementRates    map[data.Platform]float64 `json:"engagement_rates"`
	GrowthRates        map[data.Platform]float64 `json:"growth_rates"`
	Recommendations    []string                  `json:"recommendations"`
}
