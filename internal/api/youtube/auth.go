// Package youtube provides authentication utilities for the YouTube API.
package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// OAuth2Token represents an OAuth 2.0 token for YouTube API access.
type OAuth2Token struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// IsExpired checks if the token has expired (with a 5-minute buffer).
func (t *OAuth2Token) IsExpired() bool {
	return time.Now().Add(5 * time.Minute).After(t.ExpiresAt)
}

// Authenticator handles OAuth 2.0 authentication for YouTube API.
type Authenticator struct {
	clientID     string
	clientSecret string
	redirectURI  string
	token        *OAuth2Token
}

// NewAuthenticator creates a new YouTube authenticator.
func NewAuthenticator(clientID, clientSecret, redirectURI string) *Authenticator {
	return &Authenticator{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
	}
}

// GetAuthorizationURL returns the URL to redirect users for OAuth consent.
// Scopes typically needed:
// - https://www.googleapis.com/auth/youtube.readonly (read channel/video data)
// - https://www.googleapis.com/auth/yt-analytics.readonly (read analytics)
func (a *Authenticator) GetAuthorizationURL(state string, scopes []string) string {
	params := url.Values{
		"client_id":     {a.clientID},
		"redirect_uri":  {a.redirectURI},
		"response_type": {"code"},
		"scope":         {strings.Join(scopes, " ")},
		"access_type":   {"offline"},
		"state":         {state},
	}
	return "https://accounts.google.com/o/oauth2/v2/auth?" + params.Encode()
}

// ExchangeCode exchanges an authorization code for access and refresh tokens.
func (a *Authenticator) ExchangeCode(ctx context.Context, code string) (*OAuth2Token, error) {
	data := url.Values{
		"client_id":     {a.clientID},
		"client_secret": {a.clientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {a.redirectURI},
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://oauth2.googleapis.com/token",
		strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("creating token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("exchanging code: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed with status: %d", resp.StatusCode)
	}

	var token OAuth2Token
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("decoding token response: %w", err)
	}

	token.ExpiresAt = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
	a.token = &token
	return &token, nil
}

// RefreshAccessToken refreshes the access token using the refresh token.
func (a *Authenticator) RefreshAccessToken(ctx context.Context, refreshToken string) (*OAuth2Token, error) {
	data := url.Values{
		"client_id":     {a.clientID},
		"client_secret": {a.clientSecret},
		"refresh_token": {refreshToken},
		"grant_type":    {"refresh_token"},
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://oauth2.googleapis.com/token",
		strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("creating refresh request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("refreshing token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token refresh failed with status: %d", resp.StatusCode)
	}

	var token OAuth2Token
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("decoding refresh response: %w", err)
	}

	token.ExpiresAt = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
	token.RefreshToken = refreshToken // Preserve refresh token
	a.token = &token
	return &token, nil
}

// GetToken returns the current token, refreshing if necessary.
func (a *Authenticator) GetToken(ctx context.Context) (*OAuth2Token, error) {
	if a.token == nil {
		return nil, fmt.Errorf("no token available")
	}
	if a.token.IsExpired() && a.token.RefreshToken != "" {
		return a.RefreshAccessToken(ctx, a.token.RefreshToken)
	}
	return a.token, nil
}

// SetToken sets the current token (useful for loading from storage).
func (a *Authenticator) SetToken(token *OAuth2Token) {
	a.token = token
}
