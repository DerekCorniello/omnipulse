// Package handlers provides HTTP handlers for AI insights.
package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/omnipulse/omnipulse/internal/data"
	"github.com/omnipulse/omnipulse/internal/insights"
	"github.com/omnipulse/omnipulse/internal/storage"
)

// InsightsHandler handles AI insights-related HTTP requests.
type InsightsHandler struct {
	store     storage.Store
	llm       *insights.LLMClient
	templates *template.Template
}

// NewInsightsHandler creates a new InsightsHandler.
func NewInsightsHandler(store storage.Store, llm *insights.LLMClient, templates *template.Template) *InsightsHandler {
	return &InsightsHandler{
		store:     store,
		llm:       llm,
		templates: templates,
	}
}

// Index serves the insights page.
func (h *InsightsHandler) Index(w http.ResponseWriter, r *http.Request) {
	isHTMX := r.Header.Get("HX-Request") == "true"

	insightsList, err := h.store.GetRecentInsights(r.Context(), 20)
	if err != nil {
		log.Printf("error getting insights: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	templateName := "base"
	if isHTMX {
		templateName = "insights"
	}

	if err := h.templates.ExecuteTemplate(w, templateName, map[string]interface{}{
		"Page":     "insights",
		"Title":    "AI Insights",
		"Insights": insightsList,
	}); err != nil {
		log.Printf("error rendering template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// Generate triggers generation of new insights.
func (h *InsightsHandler) Generate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get current analytics summary
	dateRange := data.DateRange{} // TODO: Set appropriate date range
	summary, err := h.store.GetAnalyticsSummary(r.Context(), dateRange)
	if err != nil {
		log.Printf("error getting analytics summary: %v", err)
		http.Error(w, "Failed to get analytics data", http.StatusInternalServerError)
		return
	}

	// Generate new insight
	insight, err := h.llm.GenerateAnalyticsInsight(r.Context(), summary)
	if err != nil {
		log.Printf("error generating insight: %v", err)
		http.Error(w, "Failed to generate insight", http.StatusInternalServerError)
		return
	}

	// Save the insight
	if err := h.store.SaveInsight(r.Context(), insight); err != nil {
		log.Printf("error saving insight: %v", err)
		// Continue - we can still return the insight even if save fails
	}

	// Return the new insight as HTML
	if err := h.templates.ExecuteTemplate(w, "insight_card", insight); err != nil {
		log.Printf("error rendering template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// List handles HTMX requests for the insights list.
func (h *InsightsHandler) List(w http.ResponseWriter, r *http.Request) {
	platformStr := r.URL.Query().Get("platform")
	limit := 10 // Default

	var insightsList []*data.Insight
	var err error

	if platformStr != "" {
		insightsList, err = h.store.GetInsights(r.Context(), data.Platform(platformStr), limit)
	} else {
		insightsList, err = h.store.GetRecentInsights(r.Context(), limit)
	}

	if err != nil {
		log.Printf("error getting insights: %v", err)
		http.Error(w, "Failed to get insights", http.StatusInternalServerError)
		return
	}

	if err := h.templates.ExecuteTemplate(w, "insights_list", insightsList); err != nil {
		log.Printf("error rendering template: %v", err)
	}
}

// Suggestions handles requests for content suggestions.
func (h *InsightsHandler) Suggestions(w http.ResponseWriter, r *http.Request) {
	platformStr := r.URL.Query().Get("platform")
	if platformStr == "" {
		http.Error(w, "Platform required", http.StatusBadRequest)
		return
	}

	// TODO: Get trend data and generate suggestions
	suggestions, err := h.llm.GenerateContentSuggestions(r.Context(), data.Platform(platformStr), nil)
	if err != nil {
		log.Printf("error generating suggestions: %v", err)
		http.Error(w, "Failed to generate suggestions", http.StatusInternalServerError)
		return
	}

	if err := h.templates.ExecuteTemplate(w, "suggestions_list", suggestions); err != nil {
		log.Printf("error rendering template: %v", err)
	}
}
