// Package handlers provides HTTP handlers for platform-specific views.
package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/crossforge/crossforge/internal/data"
	"github.com/crossforge/crossforge/internal/insights"
	"github.com/crossforge/crossforge/internal/storage"
)

// PlatformHandler handles platform-specific HTTP requests.
type PlatformHandler struct {
	store      storage.Store
	aggregator *insights.Aggregator
	templates  *template.Template
}

// NewPlatformHandler creates a new PlatformHandler.
func NewPlatformHandler(store storage.Store, aggregator *insights.Aggregator, templates *template.Template) *PlatformHandler {
	return &PlatformHandler{
		store:      store,
		aggregator: aggregator,
		templates:  templates,
	}
}

// YouTube handles requests for the YouTube analytics page.
func (h *PlatformHandler) YouTube(w http.ResponseWriter, r *http.Request) {
	h.renderPlatformPage(w, r, data.PlatformYouTube, "YouTube Analytics")
}

// X handles requests for the X (Twitter) analytics page.
func (h *PlatformHandler) X(w http.ResponseWriter, r *http.Request) {
	h.renderPlatformPage(w, r, data.PlatformX, "X Analytics")
}

// LinkedIn handles requests for the LinkedIn analytics page.
func (h *PlatformHandler) LinkedIn(w http.ResponseWriter, r *http.Request) {
	h.renderPlatformPage(w, r, data.PlatformLinkedIn, "LinkedIn Analytics")
}

// renderPlatformPage renders a platform-specific analytics page.
func (h *PlatformHandler) renderPlatformPage(w http.ResponseWriter, r *http.Request, platform data.Platform, title string) {
	isHTMX := r.Header.Get("HX-Request") == "true"

	analytics, err := h.aggregator.GetPlatformAnalytics(r.Context(), platform, 30)
	if err != nil {
		log.Printf("error getting %s analytics: %v", platform, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	templateName := "base"
	if isHTMX {
		templateName = "platform"
	}

	if err := h.templates.ExecuteTemplate(w, templateName, map[string]interface{}{
		"Page":      string(platform),
		"Title":     title,
		"Platform":  platform,
		"Analytics": analytics,
	}); err != nil {
		log.Printf("error rendering template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// PlatformCard handles HTMX requests for individual platform cards.
func (h *PlatformHandler) PlatformCard(w http.ResponseWriter, r *http.Request) {
	platformStr := r.URL.Query().Get("platform")
	if platformStr == "" {
		http.Error(w, "Platform required", http.StatusBadRequest)
		return
	}

	platform := data.Platform(platformStr)
	analytics, err := h.aggregator.GetPlatformAnalytics(r.Context(), platform, 7)
	if err != nil {
		log.Printf("error getting platform card data: %v", err)
		http.Error(w, "Failed to get platform data", http.StatusInternalServerError)
		return
	}

	if err := h.templates.ExecuteTemplate(w, "platform_card", map[string]interface{}{
		"Platform":  platform,
		"Analytics": analytics,
	}); err != nil {
		log.Printf("error rendering template: %v", err)
	}
}

// Content handles HTMX requests for platform content lists.
func (h *PlatformHandler) Content(w http.ResponseWriter, r *http.Request) {
	platformStr := r.URL.Query().Get("platform")
	if platformStr == "" {
		http.Error(w, "Platform required", http.StatusBadRequest)
		return
	}

	platform := data.Platform(platformStr)
	limit := 10  // Default limit
	offset := 0  // Default offset

	// TODO: Parse limit and offset from query parameters

	var content interface{}
	var err error

	switch platform {
	case data.PlatformYouTube:
		content, err = h.store.GetVideos(r.Context(), limit, offset)
	case data.PlatformX:
		content, err = h.store.GetTweets(r.Context(), limit, offset)
	case data.PlatformLinkedIn:
		content, err = h.store.GetLinkedInPosts(r.Context(), limit, offset)
	default:
		http.Error(w, "Invalid platform", http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Printf("error getting content: %v", err)
		http.Error(w, "Failed to get content", http.StatusInternalServerError)
		return
	}

	if err := h.templates.ExecuteTemplate(w, "content_list", map[string]interface{}{
		"Platform": platform,
		"Content":  content,
	}); err != nil {
		log.Printf("error rendering template: %v", err)
	}
}
