// Package config handles application configuration for OmniPulse.
package config

import (
	"fmt"
	"os"
	"time"
)

// Config holds all application configuration settings.
type Config struct {
	// Server settings
	Server ServerConfig

	// Database settings
	Database DatabaseConfig

	// Platform API configurations
	YouTube  YouTubeConfig
	X        XConfig
	LinkedIn LinkedInConfig

	// LLM settings
	LLM LLMConfig

	// Scheduler settings
	Scheduler SchedulerConfig
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	Host string
	Port int
}

// DatabaseConfig holds SQLite database configuration.
type DatabaseConfig struct {
	Path string
}

// YouTubeConfig holds YouTube API configuration.
type YouTubeConfig struct {
	APIKey       string
	ClientID     string
	ClientSecret string
	RefreshToken string
	ChannelID    string
}

// XConfig holds X (Twitter) API configuration.
type XConfig struct {
	APIKey       string
	APISecret    string
	AccessToken  string
	AccessSecret string
	BearerToken  string
	UserID       string
}

// LinkedInConfig holds LinkedIn API configuration.
type LinkedInConfig struct {
	ClientID     string
	ClientSecret string
	AccessToken  string
	RefreshToken string
	PersonURN    string
}

// LLMConfig holds LLM (Ollama) configuration.
type LLMConfig struct {
	Endpoint string
	Model    string
	Timeout  time.Duration
}

// SchedulerConfig holds scheduler configuration.
type SchedulerConfig struct {
	FetchInterval   time.Duration
	InsightInterval time.Duration
}

// Load loads configuration from environment variables.
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnvInt("SERVER_PORT", 8080),
		},
		Database: DatabaseConfig{
			Path: getEnv("DATABASE_PATH", "./data/omnipulse.db"),
		},
		YouTube: YouTubeConfig{
			APIKey:       os.Getenv("YOUTUBE_API_KEY"),
			ClientID:     os.Getenv("YOUTUBE_CLIENT_ID"),
			ClientSecret: os.Getenv("YOUTUBE_CLIENT_SECRET"),
			RefreshToken: os.Getenv("YOUTUBE_REFRESH_TOKEN"),
			ChannelID:    os.Getenv("YOUTUBE_CHANNEL_ID"),
		},
		X: XConfig{
			APIKey:       os.Getenv("X_API_KEY"),
			APISecret:    os.Getenv("X_API_SECRET"),
			AccessToken:  os.Getenv("X_ACCESS_TOKEN"),
			AccessSecret: os.Getenv("X_ACCESS_SECRET"),
			BearerToken:  os.Getenv("X_BEARER_TOKEN"),
			UserID:       os.Getenv("X_USER_ID"),
		},
		LinkedIn: LinkedInConfig{
			ClientID:     os.Getenv("LINKEDIN_CLIENT_ID"),
			ClientSecret: os.Getenv("LINKEDIN_CLIENT_SECRET"),
			AccessToken:  os.Getenv("LINKEDIN_ACCESS_TOKEN"),
			RefreshToken: os.Getenv("LINKEDIN_REFRESH_TOKEN"),
			PersonURN:    os.Getenv("LINKEDIN_PERSON_URN"),
		},
		LLM: LLMConfig{
			Endpoint: getEnv("LLM_ENDPOINT", "http://localhost:11434"),
			Model:    getEnv("LLM_MODEL", "llama3"),
			Timeout:  time.Duration(getEnvInt("LLM_TIMEOUT_SECONDS", 30)) * time.Second,
		},
		Scheduler: SchedulerConfig{
			FetchInterval:   time.Duration(getEnvInt("FETCH_INTERVAL_MINUTES", 60)) * time.Minute,
			InsightInterval: time.Duration(getEnvInt("INSIGHT_INTERVAL_HOURS", 24)) * time.Hour,
		},
	}

	return cfg, nil
}

// getEnv returns the value of an environment variable or a default value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt returns the integer value of an environment variable or a default value.
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
			return result
		}
	}
	return defaultValue
}
