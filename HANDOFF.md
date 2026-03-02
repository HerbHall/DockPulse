# DockPulse — Claude Code Handoff

## What This Is

DockPulse is a Docker Desktop extension that checks if your running container images have newer versions available on their registries. It fills a real gap — Docker Desktop shows container status but has zero native update-checking functionality.

**Owner**: Herb Hall (github.com/HerbHall)
**License**: MIT
**Status**: Pre-development — scaffold only, no code yet
**Build Order**: Second project (after RunNotes, before PacketDeck)

---

## What Already Exists (D:\devspace\DockPulse)

- ✅ `CLAUDE.md` — Project context for Claude Code
- ✅ `HANDOFF.md` — This file
- ✅ `README.md` — Project overview
- ✅ `CONTRIBUTING.md` — Contribution guidelines
- ✅ `CHANGELOG.md` — Keep-a-Changelog format
- ✅ `LICENSE` — MIT
- ✅ `VERSION` — 0.1.0
- ✅ `.gitignore` — Comprehensive
- ✅ `.editorconfig` — Workspace standard
- ✅ `metadata.json` — Docker extension metadata
- ✅ `Dockerfile` — Extension image stub (labels set, stages TODO)
- ✅ `Makefile` — Build/install/push targets
- ✅ `docker.svg` — Placeholder icon
- ✅ `docs/FEASIBILITY.md` — Full feasibility assessment

---

## What Still Needs To Be Done

### 1. GitHub Repository (FIRST PRIORITY)

```powershell
cd /d D:\devspace\DockPulse
cmd /c "gh repo create HerbHall/DockPulse --public --source=. --remote=origin --description "Docker Desktop extension — check if your container images have updates available""
git add -A
git commit -m "chore: initial project scaffold"
git push -u origin main
```

### 2. Create GitHub Issues

Suggested issue backlog (create with `gh issue create`):

1. **Data model: registry check results** (mvp) — Define schema for storing check results, timestamps, digest comparisons
2. **React UI: container list with update status** (mvp) — Container table showing image, current digest, latest digest, status badge (up-to-date/update-available/unknown)
3. **Backend service: registry API client** (mvp) — Query Docker Hub v2 API for manifest digests, compare against running container image IDs
4. **Docker Hub authentication and rate limiting** (mvp) — Handle anonymous rate limits (100 pulls/6hrs), optional token auth for higher limits
5. **GHCR and private registry support** (enhancement) — Extend registry client beyond Docker Hub to GitHub Container Registry, self-hosted registries
6. **Scheduled background checks** (enhancement) — Periodic update checking with configurable interval
7. **Notification system** (enhancement) — Badge/alert when updates found, optional desktop notifications
8. **Update history log** (enhancement) — Track when updates were detected, when containers were updated
9. **One-click update action** (enhancement) — Pull new image + recreate container (with safety confirmation)
10. **Docker Hub publishing** (chore) — Multi-arch build, marketplace listing, screenshots

### 3. Source Directories (Create When Development Begins)

```text
ui/           — React frontend (create when starting issue #2)
backend/      — Go or Node backend service (create when starting issue #3)
```

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

## Suggested First Session Plan

1. Create GitHub repo + push scaffold
2. Create labels + issues
3. Begin MVP: registry API client (issue #3) — this is the core technical challenge
4. Build UI showing container update status (issue #2)
5. Wire frontend to backend (issue #1)
