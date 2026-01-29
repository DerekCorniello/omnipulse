// Package handlers provides HTTP handlers for the HTMX frontend.
package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/crossforge/crossforge/internal/insights"
	"github.com/crossforge/crossforge/internal/storage"
)

// DashboardHandler handles dashboard-related HTTP requests.
type DashboardHandler struct {
	store      storage.Store
	aggregator *insights.Aggregator
	templates  *template.Template
}

// NewDashboardHandler creates a new DashboardHandler.
func NewDashboardHandler(store storage.Store, aggregator *insights.Aggregator, templates *template.Template) *DashboardHandler {
	return &DashboardHandler{
		store:      store,
		aggregator: aggregator,
		templates:  templates,
	}
}

// Index serves the main dashboard page.
func (h *DashboardHandler) Index(w http.ResponseWriter, r *http.Request) {
	// For full page loads, render the complete HTML
	// For HTMX requests, only render the content partial
	isHTMX := r.Header.Get("HX-Request") == "true"

	data, err := h.aggregator.GetDashboardData(r.Context(), 7) // Last 7 days
	if err != nil {
		log.Printf("error getting dashboard data: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	templateName := "base"
	if isHTMX {
		templateName = "dashboard"
	}

	if err := h.templates.ExecuteTemplate(w, templateName, map[string]interface{}{
		"Page":       "dashboard",
		"Title":      "Dashboard",
		"Data":       data,
	}); err != nil {
		log.Printf("error rendering template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// Refresh handles HTMX requests to refresh dashboard data.
func (h *DashboardHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	days := 7
	if d := r.URL.Query().Get("days"); d != "" {
		// TODO: Parse days parameter
	}

	data, err := h.aggregator.GetDashboardData(r.Context(), days)
	if err != nil {
		log.Printf("error getting dashboard data: %v", err)
		http.Error(w, "Failed to refresh data", http.StatusInternalServerError)
		return
	}

	if err := h.templates.ExecuteTemplate(w, "dashboard_content", data); err != nil {
		log.Printf("error rendering template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// Summary handles requests for the analytics summary card.
func (h *DashboardHandler) Summary(w http.ResponseWriter, r *http.Request) {
	data, err := h.aggregator.GetDashboardData(r.Context(), 7)
	if err != nil {
		log.Printf("error getting summary: %v", err)
		http.Error(w, "Failed to get summary", http.StatusInternalServerError)
		return
	}

	if err := h.templates.ExecuteTemplate(w, "summary_card", data.Summary); err != nil {
		log.Printf("error rendering template: %v", err)
	}
}
