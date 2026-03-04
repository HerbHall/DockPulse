# DockPulse -- Copilot Instructions

Docker Desktop extension for monitoring Docker container health and performance. Queries Docker Hub v2 API directly for image digest comparison to detect outdated container images.

## Tech Stack

- **Backend**: Go 1.24, stdlib `net/http`, Docker Engine API (`github.com/docker/docker`), SQLite (`modernc.org/sqlite`)
- **Frontend**: React 18, TypeScript 5, MUI v5 (pinned via `@docker/docker-mui-theme`), Vite 7, Vitest
- **Communication**: Unix socket via `@docker/extension-api-client`
- **Build**: Multi-stage Dockerfile, Makefile, GitHub Actions CI

## Project Structure

```text
DockPulse/
├── backend/
│   ├── main.go                       - Server entry point (Unix socket listener)
│   └── internal/
│       ├── api/handler.go            - HTTP handlers (3 endpoints)
│       ├── checker/checker.go        - Orchestrator (enumerate, compare, store)
│       ├── docker/client.go          - Docker Engine API wrapper
│       ├── imageref/parse.go         - Image reference parser
│       ├── registry/                 - Registry interface + Docker Hub v2 client
│       └── store/                    - SQLite store (models, queries, migrations)
├── ui/
│   └── src/
│       ├── App.tsx                   - Main layout
│       ├── hooks/useBackend.ts       - ddClient API calls
│       ├── types.ts                  - Shared types
│       └── components/               - ContainerTable, StatusChip, ErrorBoundary
├── Dockerfile                        - Multi-stage build
├── Makefile                          - Build/test/lint targets
└── metadata.json                     - Docker extension metadata
```

## Code Style

- Conventional commits: `feat:`, `fix:`, `refactor:`, `docs:`, `test:`, `chore:`
- Co-author tag: `Co-Authored-By: GitHub Copilot <noreply@github.com>`
- Go errors wrapped with context: `fmt.Errorf("operation: %w", err)`
- Table-driven tests with `t.Run` and descriptive names
- All lint checks must pass before committing (golangci-lint v2 for Go, ESLint for TS)
- Backend Go commands use `cd backend` prefix (go.mod is in `backend/`, not root)

## Coding Guidelines

- Fix errors immediately -- never classify them as pre-existing
- Build, test, and lint must pass before any commit
- Never skip hooks (`--no-verify`) or force-push main
- Validate only at system boundaries (user input, external APIs)
- Remove unused code completely; no backwards-compatibility hacks
- MUI v5 only: use `InputProps` (not `slotProps.input`) for TextField adornments
- Docker extension UI reinitializes on every tab switch -- all state must come from backend

## Available Resources

```bash
make validate         # Run all checks (build, test, lint, typecheck)
make go-build         # cd backend && go build ./...
make go-test          # cd backend && go test ./...
make go-lint          # cd backend && golangci-lint v2 run ./...
make fe-build         # cd ui && npm run build
make fe-test          # cd ui && npx vitest run
make fe-lint          # cd ui && npx eslint src/
make fe-typecheck     # cd ui && npx tsc --noEmit
make build-extension  # Build Docker extension image
```

## Do NOT

- Add `//nolint` directives without fixing the root cause first
- Use `any` in TypeScript or suppress TypeScript errors with `as unknown`
- Commit generated files without regenerating them first
- Add dependencies without updating the lock file
- Use `panic` in library code; return errors instead
- Store secrets, tokens, or credentials in code or config files
- Mark work as complete when known errors remain
- Run Go commands from the repo root -- always `cd backend` first
