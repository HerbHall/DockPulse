<!--
  Scope: AGENTS.md guides the Copilot coding agent and Copilot Chat.
  For code completion and code review patterns, see .github/copilot-instructions.md
  and .github/instructions/*.instructions.md
  For Claude Code, see CLAUDE.md
-->

# DockPulse

Docker Desktop extension for monitoring Docker container health and performance. Queries Docker Hub v2 API directly for image digest comparison -- detects when running containers have outdated images.

## Tech Stack

- **Backend**: Go 1.24, stdlib `net/http`, Docker Engine API (`github.com/docker/docker`), SQLite (`modernc.org/sqlite`)
- **Frontend**: React 18, TypeScript 5, MUI v5 (pinned via `@docker/docker-mui-theme`), Vite 7, Vitest
- **Communication**: Unix socket via `@docker/extension-api-client` (`ddClient.extension.vm.service.get/post()`)
- **Build**: Multi-stage Dockerfile, Makefile, GitHub Actions CI (7 parallel jobs)
- **Platform**: Docker Desktop 4.8.0+ (Windows, macOS, Linux)

## Build and Test Commands

```bash
# Full verification (run before any PR)
make validate            # Build + test + lint + typecheck (all targets)

# Backend
make go-build            # cd backend && go build ./...
make go-test             # cd backend && go test ./...
make go-lint             # cd backend && golangci-lint v2 run ./...

# Frontend
make fe-build            # cd ui && npm run build
make fe-test             # cd ui && npx vitest run
make fe-lint             # cd ui && npx eslint src/
make fe-typecheck        # cd ui && npx tsc --noEmit

# Docker extension
make build-extension     # Build Docker extension image
make install-extension   # Install into Docker Desktop
make push-extension      # Multi-arch build and push (linux/amd64 + linux/arm64)
```

## Project Structure

```text
DockPulse/
├── backend/
│   ├── main.go                       - Server entry point (Unix socket listener)
│   └── internal/
│       ├── api/handler.go            - HTTP handlers (3 endpoints)
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

## Workflow Rules

### Always Do

- Create a feature branch for every change (`feature/issue-NNN-description`)
- Use conventional commits: `feat:`, `fix:`, `refactor:`, `docs:`, `test:`, `chore:`
- Run `make validate` before opening a PR
- Write table-driven tests with descriptive names (Go)
- Wrap errors with context: `fmt.Errorf("operation: %w", err)`
- Fix every error you find, regardless of who introduced it
- Backend Go code lives in `backend/` -- always `cd backend` before Go commands

### Ask First

- Adding new dependencies (check if stdlib covers the need)
- Architectural changes (new packages, major interface changes)
- Database schema migrations
- Changes to CI/CD workflows
- Removing or renaming public APIs
- Adding new API endpoints

### Never Do

- Commit directly to `main` -- always use feature branches
- Skip tests or lint checks -- even for "small changes"
- Use `--no-verify` or `--force` flags
- Commit secrets, credentials, or API keys
- Add TODO comments without a linked issue number
- Mark work as complete when build, test, or lint failures remain
- Use `panic` in library code; return errors instead
- Use `any` in TypeScript; use `unknown` with type guards

## Core Principles

These are unconditional -- no optimization or time pressure overrides them:

1. **Quality**: Once found, always fix, never leave. There is no "pre-existing" error.
2. **Verification**: Build, test, and lint must pass before any commit.
3. **Safety**: Never force-push `main`. Never skip hooks. Never commit secrets.
4. **Honesty**: Never mark work as complete when it is not.

## Error Handling

```go
// Wrap errors with context -- every return site should add meaning
if err != nil {
    return fmt.Errorf("load config: %w", err)
}

// Use sentinel errors for caller-distinguishable conditions
var ErrNotFound = errors.New("not found")
if errors.Is(err, ErrNotFound) { ... }
```

## Testing Conventions

```go
// Table-driven tests with descriptive names
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {
            name:  "valid input returns expected output",
            input: "example",
            want:  "result",
        },
        {
            name:    "empty input returns error",
            input:   "",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := FunctionName(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("FunctionName() error = %v, wantErr %v", err, tt.wantErr)
            }
            if got != tt.want {
                t.Errorf("FunctionName() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## React/TypeScript Conventions

```tsx
// MUI v5 (pinned via @docker/docker-mui-theme) -- use InputProps, NOT slotProps
<TextField InputProps={{ startAdornment: <InputAdornment position="start">...</InputAdornment> }} />

// Docker Desktop extension API -- use ddClient, never raw fetch
import { createDockerDesktopClient } from '@docker/extension-api-client'
const ddClient = createDockerDesktopClient()
const result = await ddClient.extension.vm.service.get('/api/checks')
```

## Commit Format

```text
feat: add container health monitoring endpoint

Implements periodic health checks with configurable intervals.

Closes #42
Co-Authored-By: GitHub Copilot <copilot@github.com>
```

Types: `feat` (new feature), `fix` (bug fix), `refactor` (no behavior change),
`docs` (documentation only), `test` (tests only), `chore` (build/tooling).
