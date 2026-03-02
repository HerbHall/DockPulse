# DockPulse

Docker Desktop extension — container image update checker.

## Tech Stack

- **Frontend**: React 18 + MUI v5 (pinned via `@docker/docker-mui-theme`) + Vite 7 + TypeScript 5
- **Backend**: Go 1.24, stdlib `net/http`, Docker Engine API client (`github.com/docker/docker`)
- **Storage**: SQLite via `modernc.org/sqlite` (pure Go, no CGO) at `/data/dockpulse.db`
- **Registry**: Docker Hub v2 API — anonymous token auth, manifest HEAD for digest comparison
- **Communication**: Unix socket via `ddClient.extension.vm.service.get/post()`
- **Build**: Multi-stage Dockerfile + Makefile + GitHub Actions CI (7 parallel jobs)
- **Platform**: Docker Desktop 4.8.0+ (Windows, Mac, Linux)

## Key Design Decisions

- Query registries directly for digest comparison — do NOT wrap Watchtower or any external tool
- Docker Hub v2 API for manifest digest retrieval
- Compare running image RepoDigests against registry latest
- Storage lives in a Docker volume attached to the backend container
- Extension UI reinitializes on every tab switch — all state must come from backend

## Project Conventions

- Commit messages: conventional commits (`feat:`, `fix:`, `docs:`, `chore:`)
- Co-authored commits with Claude: `Co-Authored-By: Claude <noreply@anthropic.com>`
- Issues track all work; PRs reference issue numbers
- PowerShell is the primary scripting shell on Windows

## Build Commands

```bash
make validate            # Run all checks (build, test, lint, typecheck)
make build-extension     # Build the Docker extension image
make install-extension   # Install into Docker Desktop
make go-test             # Backend tests
make go-lint             # Backend lint (golangci-lint v2)
make fe-test             # Frontend tests (vitest)
make fe-typecheck        # TypeScript type checking
make fe-lint             # Frontend lint (eslint)
```

## API Endpoints

| Method | Path | Purpose |
| -------- | ------ | --------- |
| GET | `/api/checks` | Latest check results for all containers |
| POST | `/api/check-all` | Trigger check for all running containers |
| GET | `/api/status` | Health check |

## File Layout

```text
DockPulse/
├── backend/
│   ├── main.go                       - Server entry point (Unix socket listener)
│   └── internal/
│       ├── api/handler.go            - HTTP handlers for 3 endpoints
│       ├── checker/checker.go        - Orchestrator (enumerate, compare, store)
│       ├── docker/client.go          - Docker Engine API wrapper
│       ├── imageref/parse.go         - Image reference parser (registry/ns/name:tag)
│       ├── registry/registry.go      - Registry interface
│       ├── registry/dockerhub.go     - Docker Hub v2 client (token auth + manifest HEAD)
│       └── store/                    - SQLite store (models, queries, migrations)
├── ui/
│   └── src/
│       ├── App.tsx                   - Main layout with header, table, check button
│       ├── hooks/useBackend.ts       - ddClient API calls (GET/POST to backend)
│       ├── types.ts                  - Shared TypeScript types
│       └── components/               - ContainerTable, StatusChip, ErrorBoundary
├── .github/workflows/ci.yml         - 7-job CI pipeline
├── metadata.json                     - Docker extension metadata
├── Dockerfile                        - Multi-stage build (Go + Node + Alpine)
├── Makefile                          - Build/test/lint targets
└── docs/FEASIBILITY.md              - Original feasibility research
```
