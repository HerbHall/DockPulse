# DockPulse — Claude Code Handoff

## What This Is

DockPulse is a Docker Desktop extension that checks if your running container images have newer versions available on their registries. It fills a real gap — Docker Desktop shows container status but has zero native update-checking functionality.

**Owner**: Herb Hall (github.com/HerbHall)
**License**: MIT
**Status**: Phase 1 MVP complete — backend, frontend, CI, Dockerfile all working
**Build Order**: Second project (after RunNotes, before PacketDeck)

---

## What Already Exists (D:\DevSpace\DockPulse)

### Scaffold (pre-development)

- ✅ `CLAUDE.md`, `HANDOFF.md`, `README.md`, `CONTRIBUTING.md`
- ✅ `CHANGELOG.md`, `LICENSE` (MIT), `VERSION` (0.1.0)
- ✅ `.gitignore`, `.editorconfig`, `metadata.json`
- ✅ `docs/FEASIBILITY.md` — Full feasibility assessment
- ✅ `docker.svg` — Placeholder icon

### Phase 1 MVP (completed)

- ✅ `.github/workflows/ci.yml` — 7-job CI pipeline (go-lint, go-build, go-test, fe-lint, fe-typecheck, fe-test, docker-build)
- ✅ `backend/` — Go backend with 15 source files across 6 packages
  - `internal/store/` — SQLite data model with image check tracking
  - `internal/imageref/` — Image reference parser (Docker Hub, GHCR, private registries)
  - `internal/registry/` — Docker Hub v2 client (token auth, manifest HEAD, digest extraction)
  - `internal/docker/` — Docker Engine API wrapper (container list, image inspect)
  - `internal/checker/` — Orchestrator tying store, registry, and Docker client together
  - `internal/api/` — HTTP handlers (GET /api/checks, POST /api/check-all, GET /api/status)
  - `main.go` — Server entry point with Unix socket listener and graceful shutdown
- ✅ `ui/` — React 18 + MUI v5 + Vite 7 + TypeScript frontend
  - Components: ContainerTable, StatusChip, ErrorBoundary
  - Hook: useBackend (real ddClient API calls)
  - Tests: 5 unit tests with vitest + testing-library
- ✅ `Dockerfile` — Multi-stage build (Go 1.24-alpine → Node 22-alpine → Alpine 3.19)
- ✅ `Makefile` — Full targets including validate, go-test, go-lint, fe-test, fe-typecheck, fe-lint

---

## What Still Needs To Be Done

### Phase 1 (complete)

- ✅ GitHub repo created, issues filed (#1-#10)
- ✅ #1 Data model — SQLite store (PR #12)
- ✅ #2 React UI — container list with status (PR #14, #16)
- ✅ #3 Backend service — registry client + checker + API (PR #15)
- ✅ #4 Image ref parsing — Docker Hub, GHCR, private registry formats (PR #13)
- ✅ #5 Docker Hub auth — anonymous token auth for manifest HEAD (PR #15)
- ✅ #10 Dockerfile + branding (PR #17)

### Phase 2+ (open issues)

- #6 Scheduled background checks with configurable intervals
- #7 Scan history and update timeline
- #8 Multi-registry support (GHCR, Quay, private)
- #9 Docker Hub publishing and marketplace submission

---

## Key Architecture Decisions

These are settled from the research phase:

- **Native registry checking** — Do NOT wrap Watchtower or any external tool. Query registries directly.
- **Docker Hub v2 API** for digest comparison: `GET /v2/{repo}/manifests/{tag}` returns digest headers
- **Compare digests** between running container's image ID and registry's latest manifest digest
- **Storage in Docker volume** attached to backend container (same pattern as RunNotes)
- **React + Material UI** frontend matching Docker Desktop look
- **Socket communication** between frontend and backend (Extensions SDK standard)
- **Multi-arch required**: linux/amd64 + linux/arm64

## Registry API Key Facts

- Docker Hub v2 API: `registry-1.docker.io/v2/{namespace}/{repo}/manifests/{tag}`
- Auth flow: GET token from `auth.docker.io/token` with `scope=repository:{repo}:pull`
- Anonymous rate limit: 100 pulls per 6 hours per IP
- Authenticated rate limit: 200 pulls per 6 hours
- GHCR: `ghcr.io/v2/{owner}/{repo}/manifests/{tag}` (similar v2 API)
- Digest comparison: `Docker-Content-Digest` header vs running image's `RepoDigests`

## Name Research

"DockPulse" was chosen after conflict checks confirmed:
- No GitHub repository named DockPulse
- No Docker Hub image named dockpulse
- No npm package named dockpulse
- No trademark conflicts
- Conveys "checking the pulse/health of your Docker containers"

## Herb's Preferences (Important)

- **Green/earthy colors** for branding — dislikes blue
- **PowerShell is primary shell** on Windows
- **Use `cmd /c` wrapper** for `gh` commands in PowerShell
- **Conventional commits**: `feat:`, `fix:`, `docs:`, `chore:`
- **Co-authored commits**: `Co-Authored-By: Claude <noreply@anthropic.com>`
- **Executes steps immediately as he reads them** — put prerequisites BEFORE action steps

---

## Suggested Next Session Plan

1. Install the extension locally and test with real containers (`make build-extension && make install-extension`)
2. Take screenshots for marketplace submission (issue #9)
3. Implement scheduled scans (issue #6) — add settings UI + backend timer
4. Add GHCR support (issue #8) — extend Registry interface with a GHCR client
