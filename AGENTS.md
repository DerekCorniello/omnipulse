# AGENTS.md

This file provides guidance for AI coding agents (LLMs, Cursor, Aider, etc.) working on this project: a multiplatform content analytics tool (YouTube + X + LinkedIn) that pulls stats/comments/engagement, stores data, processes trends, and generates AI insights.

## Project Overview
- **Purpose**: Fetch, aggregate, and analyze creator metrics from YouTube, X (Twitter), and LinkedIn APIs. Provide summaries, alerts, and LLM-powered insights on how I can increase my presence and reach via these platforms, new ideas for content, etc.
- **Core Components**:
  - API clients for YouTube, X, and LinkedIn.
  - Data processing with structs/slices.
  - Embedded storage with sqlite.
  - Scheduling via time.Ticker/goroutines.
  - LLM integration will likely be done via my local instance of ollama.
  - HTMX frontend for lightweight dashboards.
- **Main Goal**: Keep it lightweight, reliable, single-binary deployable.

## Structure (recommended)
omnipulse/
├── cmd/
│   └── omnipulse/
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── youtube/
│   │   │   ├── client.go
│   │   │   ├── analytics.go
│   │   │   └── auth.go
│   │   ├── x/
│   │   │   ├── client.go
│   │   │   ├── metrics.go
│   │   │   └── replies.go
│   │   └── linkedin/
│   │       ├── client.go
│   │       ├── posts.go
│   │       └── analytics.go
│   ├── config/
│   │   └── config.go
│   ├── data/
│   │   ├── models.go
│   │   └── types.go
│   ├── insights/
│   │   ├── aggregator.go
│   │   ├── llm.go
│   │   └── trends.go
│   ├── storage/
│   │   ├── sqlite.go
│   │   └── storage.go
│   └── scheduler/
│       └── scheduler.go
├── pkg/
│   └── cli/
│       └── root.go
├── scripts/
│   └── setup-oauth.sh
├── .env.example
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
├── README.md
└── AGENTS.md

## Build & Run Commands
- Build: `go build -o bin/content-tool ./cmd/main.go`
- Run CLI: `go run ./cmd/main.go fetch --platform youtube` (or full args)
- Tests: `go test ./... -v`
- Lint: `golangci-lint run`
- Docker build: `docker build -t omnipulse .`
- Docker run: `docker run -v $(pwd)/data:/app/data omnipulse`

Always run `go test ./...` before commits/PRs.

## Code Style & Conventions
- Use Go 1.25+ features.
- Prefer explicit errors: `if err != nil { return err }`.
- Keep packages small and focused (e.g., `internal/api/youtube`, `pkg/insights`).
- No external deps unless necessary; vendor if possible.
- Comments: Godoc-style for exported items; inline for complex logic.
- Error wrapping: Use `fmt.Errorf("fetching %s: %w", platform, err)`.

## Testing Guidelines
- Unit test all API wrappers, parsers, and insight logic.
- Use table-driven tests for edge cases (e.g., API rate limits, empty responses).
- Mock HTTP clients with `httptest` or custom roundtrippers.

## Git & PR Workflow
- Branch naming: `feature/fetch-youtube-analytics`, `fix/linkedin-auth-refresh`, `chore/deps-update`.
- Commit messages: Conventional style (e.g., `feat: add X API comment fetching`, `fix: handle OAuth token expiry`).
- PR titles: `<type>: <short description>` (feat/fix/chore/docs).
- Always include tests/changes; request review for API/auth changes.
- Rebase main before merging.

## Boundaries & Preferences
- Do NOT add heavy frameworks (no Gin if Chi suffices; prefer stdlib where possible).
- Avoid external services unless for LLM (prefer local Ollama).
- No big-bang changes: small, incremental PRs.
- Favor simplicity over premature optimization.
- If adding new platforms/APIs, create a new package under `internal/api/<platform>`.

## LLM/Agent-Specific Instructions
- When suggesting code: Prefer idiomatic Go (no unnecessary abstractions).
- For insights: Suggest modular functions that can call external LLMs cleanly.
- If proposing changes: Reference existing structs (e.g., `type PlatformData struct { ... }`).
- Always prioritize security: Never hardcode API keys; use env vars or config.
