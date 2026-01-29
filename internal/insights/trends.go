// Package insights provides trend analysis for analytics data.
package insights

import (
	"context"
	"math"
	"time"

	"github.com/crossforge/crossforge/internal/data"
	"github.com/crossforge/crossforge/internal/storage"
)

// TrendAnalyzer analyzes trends in analytics data.
type TrendAnalyzer struct {
	store storage.Store
}

// NewTrendAnalyzer creates a new TrendAnalyzer.
func NewTrendAnalyzer(store storage.Store) *TrendAnalyzer {
	return &TrendAnalyzer{store: store}
}

// AnalyzeTrend calculates trend direction and percentage change.
func (t *TrendAnalyzer) AnalyzeTrend(ctx context.Context, platform data.Platform, metric string, days int) (*data.TrendData, error) {
	trend, err := t.store.GetTrendData(ctx, platform, metric, days)
	if err != nil {
		return nil, err
	}

	if trend == nil || len(trend.Points) < 2 {
		return nil, nil
	}

	// Calculate trend direction and change
	trend.Trend, trend.ChangePercent = calculateTrendMetrics(trend.Points)

	return trend, nil
}

// calculateTrendMetrics calculates trend direction and percentage change.
func calculateTrendMetrics(points []data.DataPoint) (string, float64) {
	if len(points) < 2 {
		return "stable", 0
	}

	// Sort points by timestamp (should already be sorted, but ensure)
	// Use first and last points for overall trend
	firstValue := points[0].Value
	lastValue := points[len(points)-1].Value

	if firstValue == 0 {
		if lastValue > 0 {
			return "up", 100
		}
		return "stable", 0
	}

	changePercent := ((lastValue - firstValue) / firstValue) * 100

	var trend string
	switch {
	case changePercent > 5:
		trend = "up"
	case changePercent < -5:
		trend = "down"
	default:
		trend = "stable"
	}

	return trend, math.Round(changePercent*100) / 100
}

// DetectAnomalies identifies unusual patterns in the data.
func (t *TrendAnalyzer) DetectAnomalies(ctx context.Context, platform data.Platform, metric string, days int) ([]*Anomaly, error) {
	trend, err := t.store.GetTrendData(ctx, platform, metric, days)
	if err != nil || trend == nil {
		return nil, err
	}

	anomalies := detectAnomaliesInPoints(trend.Points)
	return anomalies, nil
}

// Anomaly represents an unusual data point or pattern.
type Anomaly struct {
	Timestamp   time.Time `json:"timestamp"`
	Value       float64   `json:"value"`
	ExpectedMin float64   `json:"expected_min"`
	ExpectedMax float64   `json:"expected_max"`
	Type        string    `json:"type"` // "spike", "drop", "unusual_pattern"
	Severity    string    `json:"severity"` // "low", "medium", "high"
}

// detectAnomaliesInPoints uses simple statistical analysis to find anomalies.
func detectAnomaliesInPoints(points []data.DataPoint) []*Anomaly {
	if len(points) < 7 { // Need enough data for meaningful analysis
		return nil
	}

	// Calculate mean and standard deviation
	var sum, sumSquares float64
	for _, p := range points {
		sum += p.Value
		sumSquares += p.Value * p.Value
	}
	mean := sum / float64(len(points))
	variance := (sumSquares / float64(len(points))) - (mean * mean)
	stdDev := math.Sqrt(variance)

	// Find points outside 2 standard deviations
	var anomalies []*Anomaly
	threshold := 2.0
	for _, p := range points {
		deviation := math.Abs(p.Value - mean)
		if deviation > threshold*stdDev {
			anomalyType := "spike"
			if p.Value < mean {
				anomalyType = "drop"
			}

			severity := "low"
			if deviation > 3*stdDev {
				severity = "high"
			} else if deviation > 2.5*stdDev {
				severity = "medium"
			}

			anomalies = append(anomalies, &Anomaly{
				Timestamp:   p.Timestamp,
				Value:       p.Value,
				ExpectedMin: mean - threshold*stdDev,
				ExpectedMax: mean + threshold*stdDev,
				Type:        anomalyType,
				Severity:    severity,
			})
		}
	}

	return anomalies
}

// ComparePerioods compares metrics between two time periods.
func (t *TrendAnalyzer) ComparePeriods(ctx context.Context, platform data.Platform, metric string, period1, period2 data.DateRange) (*PeriodComparison, error) {
	// TODO: Implement period-over-period comparison
	// Useful for week-over-week or month-over-month analysis
	return nil, nil
}

// PeriodComparison holds comparison data between two time periods.
type PeriodComparison struct {
	Platform     data.Platform `json:"platform"`
	Metric       string        `json:"metric"`
	Period1Value float64       `json:"period1_value"`
	Period2Value float64       `json:"period2_value"`
	Change       float64       `json:"change"`
	ChangePercent float64      `json:"change_percent"`
	Trend        string        `json:"trend"`
}
