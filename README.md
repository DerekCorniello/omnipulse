# CrossForge

**Multiplatform Content Analytics Dashboard**

CrossForge is a lightweight, single-binary content analytics tool that aggregates metrics from YouTube, X (Twitter), and LinkedIn. It features an HTMX-powered dashboard and LLM integration (via Ollama) for AI-powered insights.

## Features

- **Unified Dashboard**: View all your social media analytics in one place
- **Platform Support**: YouTube, X (Twitter), and LinkedIn integration
- **SQLite Storage**: Embedded database for historical data and trend analysis
- **HTMX Frontend**: Lightweight, interactive dashboard without heavy JavaScript frameworks
- **LLM Insights**: AI-powered analytics insights via local Ollama instance
- **Scheduled Fetching**: Automatic background data collection
- **Docker Support**: Easy deployment with Docker and docker-compose

## Quick Start

### Prerequisites

- Go 1.22+ installed
- SQLite3 (for database management)
- Ollama (optional, for AI insights)
- API credentials for the platforms you want to track

### Installation

1. Clone the repository:
```bash
git clone https://github.com/crossforge/crossforge.git
cd crossforge
```

2. Copy the environment template and configure your API credentials:
```bash
cp .env.example .env
# Edit .env with your API keys and tokens
```

3. Build the application:
```bash
go build -o bin/crossforge ./cmd/crossforge
```

4. Run database migrations:
```bash
./scripts/migrate-db.sh
```

5. Start the server:
```bash
./bin/crossforge serve
```

6. Open your browser to `http://localhost:8080`

### Docker Deployment

```bash
# Build and run with docker-compose
docker-compose up -d

# View logs
docker-compose logs -f crossforge
```

## Project Structure

```
crossforge/
├── cmd/crossforge/          # Application entry point
├── internal/
│   ├── api/                 # Platform API clients
│   │   ├── youtube/         # YouTube Data API v3 client
│   │   ├── x/               # X (Twitter) API v2 client
│   │   └── linkedin/        # LinkedIn API client
│   ├── config/              # Configuration management
│   ├── data/                # Data models and types
│   ├── insights/            # Analytics aggregation and LLM integration
│   ├── storage/             # SQLite storage layer
│   │   └── migrations/      # Database migrations
│   ├── scheduler/           # Background task scheduling
│   └── frontend/            # HTMX handlers and templates
├── pkg/cli/                 # CLI commands
├── web/                     # Static assets (CSS, JS)
├── scripts/                 # Setup and utility scripts
└── testdata/                # Test fixtures
```

## Configuration

All configuration is done via environment variables. See `.env.example` for the full list.

### Core Settings

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_HOST` | `localhost` | Server bind address |
| `SERVER_PORT` | `8080` | Server port |
| `DATABASE_PATH` | `./data/crossforge.db` | SQLite database path |
| `LLM_ENDPOINT` | `http://localhost:11434` | Ollama endpoint |
| `LLM_MODEL` | `llama3` | Ollama model for insights |

---

# Platform API Setup Guides

## YouTube API Setup

### Required APIs
- **YouTube Data API v3**: Channel/video statistics, comments
- **YouTube Analytics API**: Detailed analytics (views, watch time, subscribers)

### Getting Started

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the YouTube Data API v3 and YouTube Analytics API
4. Create credentials:
   - **API Key**: For public data access
   - **OAuth 2.0 Client ID**: For analytics and private data

### OAuth 2.0 Scopes Required

| Scope | Purpose |
|-------|---------|
| `youtube.readonly` | Read channel and video data |
| `yt-analytics.readonly` | Read analytics reports |
| `yt-analytics-monetary.readonly` | Read revenue data (optional) |

### Rate Limits

- **Default quota**: 10,000 units per day
- `channels.list`, `videos.list`: 1 unit each
- `search.list`: 100 units (avoid when possible)
- Quota resets at midnight Pacific Time

### Key Endpoints

| Endpoint | Purpose |
|----------|---------|
| `GET /youtube/v3/channels` | Channel statistics |
| `GET /youtube/v3/videos` | Video metrics |
| `GET /youtube/v3/commentThreads` | Video comments |
| `GET /youtubeanalytics/v2/reports` | Detailed analytics |

### Getting a Refresh Token

1. Configure OAuth consent screen in Google Cloud Console
2. Create OAuth 2.0 credentials with redirect URI: `http://localhost:8080/oauth/youtube/callback`
3. Build authorization URL with `access_type=offline` and `prompt=consent`
4. Exchange authorization code for tokens
5. Store refresh token securely

### Important Limitations

- Service accounts are NOT supported by YouTube APIs
- Subscriber counts are rounded to 3 significant figures
- Analytics data may be delayed 24-72 hours
- Refresh token only returned on first authorization (use `prompt=consent`)

---

## X (Twitter) API Setup

### API Tiers

| Tier | Cost | Read Limit | Post Limit |
|------|------|------------|------------|
| Free | $0 | 100/month | 500/month |
| Basic | $200/month | 15,000/month | 50,000/month |
| Pro | $5,000/month | 10M/month | Higher limits |

**Recommendation**: Basic tier minimum for meaningful analytics.

### Getting Started

1. Go to [X Developer Portal](https://developer.x.com)
2. Create a project and app
3. From "Keys and tokens", get:
   - API Key and Secret
   - Bearer Token
   - Access Token and Secret
4. Get your User ID from [tweeterid.com](https://tweeterid.com)

### Authentication Methods

| Method | Use Case |
|--------|----------|
| Bearer Token | Read-only public data |
| OAuth 2.0 with PKCE | User context actions |
| OAuth 1.0a | Legacy endpoints, media upload |

### Required Scopes

| Scope | Purpose |
|-------|---------|
| `tweet.read` | Read tweets and metrics |
| `users.read` | Read profile and follower data |
| `offline.access` | Get refresh tokens |
| `follows.read` | Read follower/following lists |

### Rate Limits (per 15-minute window)

| Endpoint | App Limit | User Limit |
|----------|-----------|------------|
| Tweet lookup | 3,500 | 5,000 |
| User timeline | 1,500 | 900 |
| Search (recent) | 450 | 180 |

### Key Endpoints

| Endpoint | Purpose |
|----------|---------|
| `GET /2/users/:id` | User profile and follower counts |
| `GET /2/users/:id/tweets` | User's tweets with metrics |
| `GET /2/tweets/:id` | Single tweet with detailed metrics |
| `GET /2/users/:id/mentions` | Mentions of the user |

### Available Metrics

| Type | Fields | Availability |
|------|--------|--------------|
| Public | retweets, replies, likes, quotes | Any tweet |
| Non-public | impressions, clicks | Own tweets only, OAuth required |
| Organic | Combined metrics | Own tweets, last 30 days |

### Important Limitations

- Free tier insufficient for real analytics
- Non-public metrics only for tweets < 30 days old
- Full engagement API requires Enterprise tier
- Historical search requires Pro tier or higher

---

## LinkedIn API Setup

### Required Products

- **Community Management API**: Post analytics, follower stats (main product)
- **Sign In with LinkedIn**: User authentication
- **Share on LinkedIn**: Posting content

### Getting Started

1. Go to [LinkedIn Developer Portal](https://developer.linkedin.com)
2. Create a new application
3. Add required products:
   - Sign In with LinkedIn using OpenID Connect
   - Share on LinkedIn
   - Request Community Management API access
4. Configure OAuth 2.0 redirect URLs

### Application Approval

Community Management API requires:
- Business email (not personal)
- Organization's legal name and address
- Website and privacy policy
- Compliance with LinkedIn API Terms

**Timeline**: Approval typically takes 2-4 weeks.

### OAuth 2.0 Scopes

| Scope | Purpose |
|-------|---------|
| `openid`, `profile`, `email` | Basic profile access |
| `r_organization_social` | Read organization posts |
| `r_member_social` | Read member posts |
| `rw_organization_admin` | Follower statistics |
| `w_member_social` | Post content |

### Rate Limits

- Application-level daily limits (varies by endpoint)
- User-level limits on actions per day
- Returns HTTP 429 when exceeded
- Best practice: Space requests throughout the day

### Key Endpoints

| Endpoint | Purpose |
|----------|---------|
| `GET /v2/me` | Current user profile |
| `GET /rest/posts` | Organization posts |
| `GET /rest/organizationalEntityShareStatistics` | Post analytics |
| `GET /rest/organizationalEntityFollowerStatistics` | Follower data |
| `GET /v2/networkSizes/{urn}` | Follower count |

### Getting Person URN

```bash
# Call the /me endpoint after OAuth authentication
GET https://api.linkedin.com/v2/me

# Response includes "id" field
# Full URN format: urn:li:person:{id}
```

### Personal Profile vs Company Page

| Feature | Personal Profile | Company Page |
|---------|-----------------|--------------|
| Full analytics | Limited (new API 2025) | Yes |
| Run ads | No | Yes |
| API support | Basic | Comprehensive |
| Engagement rate | 2.75x higher | Lower |

### Important Limitations

- Strict approval process
- Many features require Partner Program
- Personal profile analytics limited
- API versioning changes frequently
- Some demographics limited to top 100 results

---

## Build & Run Commands

```bash
# Build the binary
go build -o bin/crossforge ./cmd/crossforge

# Run CLI
go run ./cmd/crossforge fetch --platform youtube
go run ./cmd/crossforge serve

# Run tests
go test ./... -v

# Lint (requires golangci-lint)
golangci-lint run

# Docker build
docker build -t crossforge .

# Docker run
docker run -v $(pwd)/data:/app/data crossforge
```

## Development

### Code Style

- Go 1.22+ features
- Explicit error handling: `if err != nil { return err }`
- Small, focused packages
- Godoc comments for exported items
- Error wrapping: `fmt.Errorf("context: %w", err)`

### Testing

- Unit test all API wrappers and parsers
- Use table-driven tests
- Mock HTTP clients with `httptest`

### Git Workflow

- Branch naming: `feature/`, `fix/`, `chore/`
- Conventional commits: `feat:`, `fix:`, `chore:`
- Always run tests before committing
- Rebase main before merging

## License

MIT License - See LICENSE file for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Write tests for new functionality
4. Submit a pull request

---

## Resources

### YouTube
- [YouTube Data API Documentation](https://developers.google.com/youtube/v3)
- [YouTube Analytics API Documentation](https://developers.google.com/youtube/analytics)

### X (Twitter)
- [X API Documentation](https://docs.x.com)
- [X Developer Portal](https://developer.x.com)

### LinkedIn
- [LinkedIn Marketing API Documentation](https://learn.microsoft.com/en-us/linkedin/marketing)
- [LinkedIn Developer Portal](https://developer.linkedin.com)

### Ollama
- [Ollama Documentation](https://ollama.ai)
