// Package insights provides LLM integration for generating AI-powered insights.
package insights

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/omnipulse/omnipulse/internal/config"
	"github.com/omnipulse/omnipulse/internal/data"
)

// LLMClient handles communication with the LLM (Ollama) for generating insights.
type LLMClient struct {
	httpClient *http.Client
	endpoint   string
	model      string
}

// NewLLMClient creates a new LLM client.
func NewLLMClient(cfg config.LLMConfig) *LLMClient {
	return &LLMClient{
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
		endpoint: cfg.Endpoint,
		model:    cfg.Model,
	}
}

// OllamaRequest represents a request to the Ollama API.
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// OllamaResponse represents a response from the Ollama API.
type OllamaResponse struct {
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	Response  string    `json:"response"`
	Done      bool      `json:"done"`
}

// GenerateInsight generates an insight using the LLM.
func (c *LLMClient) GenerateInsight(ctx context.Context, prompt string) (string, error) {
	req := OllamaRequest{
		Model:  c.model,
		Prompt: prompt,
		Stream: false,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshaling request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST",
		c.endpoint+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("decoding response: %w", err)
	}

	return ollamaResp.Response, nil
}

// GenerateAnalyticsInsight generates an insight based on analytics data.
func (c *LLMClient) GenerateAnalyticsInsight(ctx context.Context, summary *data.AnalyticsSummary) (*data.Insight, error) {
	prompt := buildAnalyticsPrompt(summary)

	response, err := c.GenerateInsight(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("generating insight: %w", err)
	}

	insight := &data.Insight{
		ID:          generateID(),
		Platform:    "", // Cross-platform
		Type:        "summary",
		Title:       "Analytics Overview",
		Description: response,
		Confidence:  0.8,
		GeneratedAt: time.Now(),
		DataRange:   "custom",
	}

	return insight, nil
}

// GenerateContentSuggestions generates content ideas based on analytics trends.
func (c *LLMClient) GenerateContentSuggestions(ctx context.Context, platform data.Platform, trends []*data.TrendData) ([]string, error) {
	prompt := buildContentSuggestionPrompt(platform, trends)

	response, err := c.GenerateInsight(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("generating suggestions: %w", err)
	}

	// TODO: Parse response into individual suggestions
	return []string{response}, nil
}

// buildAnalyticsPrompt creates a prompt for analytics insight generation.
func buildAnalyticsPrompt(summary *data.AnalyticsSummary) string {
	// TODO: Build a detailed prompt including:
	// - Key metrics from each platform
	// - Comparison with previous periods
	// - Request for actionable insights
	return fmt.Sprintf(`Analyze the following social media analytics and provide insights:

YouTube:
- Total views: %d
- Engagement rate: %.2f%%

X (Twitter):
- Total impressions: %d
- Engagement rate: %.2f%%

LinkedIn:
- Total impressions: %d
- Engagement rate: %.2f%%

Please provide:
1. Key observations
2. Areas of strength
3. Areas needing improvement
4. Recommended actions

Keep the response concise and actionable.`,
		getYouTubeViews(summary),
		getYouTubeEngagement(summary),
		getXImpressions(summary),
		getXEngagement(summary),
		getLinkedInImpressions(summary),
		getLinkedInEngagement(summary),
	)
}

// buildContentSuggestionPrompt creates a prompt for content suggestions.
func buildContentSuggestionPrompt(platform data.Platform, trends []*data.TrendData) string {
	// TODO: Build prompt based on trending topics and performance data
	return fmt.Sprintf(`Based on recent performance trends on %s, suggest 3 content ideas that could improve engagement.`, platform)
}

// Helper functions for prompt building
func getYouTubeViews(s *data.AnalyticsSummary) int64 {
	if s != nil && s.YouTube != nil {
		return s.YouTube.TotalViews
	}
	return 0
}

func getYouTubeEngagement(s *data.AnalyticsSummary) float64 {
	if s != nil && s.YouTube != nil {
		return s.YouTube.EngagementRate * 100
	}
	return 0
}

func getXImpressions(s *data.AnalyticsSummary) int64 {
	if s != nil && s.X != nil {
		return s.X.TotalImpressions
	}
	return 0
}

func getXEngagement(s *data.AnalyticsSummary) float64 {
	if s != nil && s.X != nil {
		return s.X.EngagementRate * 100
	}
	return 0
}

func getLinkedInImpressions(s *data.AnalyticsSummary) int64 {
	if s != nil && s.LinkedIn != nil {
		return s.LinkedIn.TotalImpressions
	}
	return 0
}

func getLinkedInEngagement(s *data.AnalyticsSummary) float64 {
	if s != nil && s.LinkedIn != nil {
		return s.LinkedIn.EngagementRate * 100
	}
	return 0
}

// generateID generates a unique ID for insights.
func generateID() string {
	// TODO: Use a proper UUID library
	return fmt.Sprintf("insight_%d", time.Now().UnixNano())
}
