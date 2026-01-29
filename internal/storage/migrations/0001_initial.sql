-- CrossForge Initial Database Schema
-- Migration: 0001_initial.sql
-- Description: Create initial tables for analytics data storage

-- Enable foreign keys
PRAGMA foreign_keys = ON;

-- =============================================================================
-- YouTube Tables
-- =============================================================================

-- YouTube videos table
CREATE TABLE IF NOT EXISTS youtube_videos (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    published_at DATETIME NOT NULL,
    view_count INTEGER DEFAULT 0,
    like_count INTEGER DEFAULT 0,
    comment_count INTEGER DEFAULT 0,
    duration TEXT,
    thumbnail_url TEXT,
    fetched_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- YouTube channel stats history
CREATE TABLE IF NOT EXISTS youtube_channel_stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    subscriber_count INTEGER DEFAULT 0,
    view_count INTEGER DEFAULT 0,
    video_count INTEGER DEFAULT 0,
    fetched_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- X (Twitter) Tables
-- =============================================================================

-- X tweets table
CREATE TABLE IF NOT EXISTS x_tweets (
    id TEXT PRIMARY KEY,
    text TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    like_count INTEGER DEFAULT 0,
    retweet_count INTEGER DEFAULT 0,
    reply_count INTEGER DEFAULT 0,
    quote_count INTEGER DEFAULT 0,
    impression_count INTEGER DEFAULT 0,
    fetched_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    db_created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- X user stats history
CREATE TABLE IF NOT EXISTS x_user_stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    follower_count INTEGER DEFAULT 0,
    following_count INTEGER DEFAULT 0,
    tweet_count INTEGER DEFAULT 0,
    listed_count INTEGER DEFAULT 0,
    fetched_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- LinkedIn Tables
-- =============================================================================

-- LinkedIn posts table
CREATE TABLE IF NOT EXISTS linkedin_posts (
    id TEXT PRIMARY KEY,
    text TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    like_count INTEGER DEFAULT 0,
    comment_count INTEGER DEFAULT 0,
    share_count INTEGER DEFAULT 0,
    impression_count INTEGER DEFAULT 0,
    click_count INTEGER DEFAULT 0,
    fetched_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    db_created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- LinkedIn profile stats history
CREATE TABLE IF NOT EXISTS linkedin_profile_stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    connection_count INTEGER DEFAULT 0,
    follower_count INTEGER DEFAULT 0,
    fetched_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- Comments Table (Cross-platform)
-- =============================================================================

CREATE TABLE IF NOT EXISTS comments (
    id TEXT PRIMARY KEY,
    platform TEXT NOT NULL CHECK(platform IN ('youtube', 'x', 'linkedin')),
    content_id TEXT NOT NULL,
    author_id TEXT,
    author_name TEXT,
    text TEXT NOT NULL,
    like_count INTEGER DEFAULT 0,
    created_at DATETIME NOT NULL,
    fetched_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_comments_platform_content
ON comments(platform, content_id);

CREATE INDEX IF NOT EXISTS idx_comments_created_at
ON comments(created_at);

-- =============================================================================
-- AI Insights Table
-- =============================================================================

CREATE TABLE IF NOT EXISTS insights (
    id TEXT PRIMARY KEY,
    platform TEXT CHECK(platform IN ('youtube', 'x', 'linkedin', '')),
    type TEXT NOT NULL CHECK(type IN ('trend', 'recommendation', 'alert', 'summary')),
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    confidence REAL DEFAULT 0.0 CHECK(confidence >= 0.0 AND confidence <= 1.0),
    generated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    data_range TEXT
);

CREATE INDEX IF NOT EXISTS idx_insights_platform
ON insights(platform);

CREATE INDEX IF NOT EXISTS idx_insights_generated_at
ON insights(generated_at);

CREATE INDEX IF NOT EXISTS idx_insights_type
ON insights(type);

-- =============================================================================
-- Metrics History Table (For trend analysis)
-- =============================================================================

CREATE TABLE IF NOT EXISTS metrics_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    platform TEXT NOT NULL CHECK(platform IN ('youtube', 'x', 'linkedin')),
    metric_name TEXT NOT NULL,
    metric_value REAL NOT NULL,
    recorded_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_metrics_platform_name
ON metrics_history(platform, metric_name);

CREATE INDEX IF NOT EXISTS idx_metrics_recorded_at
ON metrics_history(recorded_at);

-- =============================================================================
-- Application Settings Table
-- =============================================================================

CREATE TABLE IF NOT EXISTS settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- OAuth Tokens Table (For storing refresh tokens securely)
-- =============================================================================

CREATE TABLE IF NOT EXISTS oauth_tokens (
    platform TEXT PRIMARY KEY CHECK(platform IN ('youtube', 'x', 'linkedin')),
    access_token TEXT NOT NULL,
    refresh_token TEXT,
    token_type TEXT,
    expires_at DATETIME,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- Schema Version Table (For migration tracking)
-- =============================================================================

CREATE TABLE IF NOT EXISTS schema_migrations (
    version INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Record this migration
INSERT INTO schema_migrations (version, name) VALUES (1, '0001_initial.sql');
