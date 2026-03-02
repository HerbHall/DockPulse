# DockPulse

Docker Desktop extension — container image update checker.

## Tech Stack

- **Frontend**: React + Material UI (Docker Desktop extension standard)
- **Backend**: TBD — Node.js or Go service running in Desktop VM
- **Storage**: SQLite or JSON file on mounted volume
- **Build**: Docker Extensions CLI + Makefile
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

## File Layout

```text
DockPulse/
├── docs/                - Feasibility research
├── ui/                  - React frontend
├── backend/             - Backend service
├── metadata.json        - Extension metadata
├── Dockerfile           - Extension image
└── Makefile             - Build targets
```
