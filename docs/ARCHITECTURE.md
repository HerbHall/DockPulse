# DockPulse Architecture

## Overview

DockPulse is a Docker Desktop extension that monitors running containers and checks whether their source images have newer versions available on the registry. It uses the standard Docker Desktop Extensions SDK pattern: React frontend tab + backend VM service + SQLite storage.

## How Update Detection Works

### Digest Comparison (Core Mechanism)

1. For each running container, extract the image reference (e.g., `nginx:latest`, `postgres:16`)
2. Query the local Docker daemon for the image's digest (SHA256)
3. Query the remote registry API for the current digest of the same tag
4. Compare: if digests differ, a newer version exists

### Registry API Flow

**Docker Hub** (v2 API):

```
GET https://auth.docker.io/token?service=registry.docker.io&scope=repository:{image}:pull
GET https://registry-1.docker.io/v2/{image}/manifests/{tag}
  → Accept: application/vnd.docker.distribution.manifest.v2+json
  → Returns: Docker-Content-Digest header
```

**OCI-compliant registries** (GHCR, ECR, ACR) follow the same v2 spec with different auth endpoints.

### What We DON'T Do

- We do NOT pull or download new images automatically
- We do NOT run Watchtower, WUD, Diun, or any external tool
- We do NOT modify containers — read-only inspection only

## Component Architecture

```
┌─────────────────────────────────────────────┐
│  Docker Desktop                             │
│  ┌───────────────────────────────────────┐  │
│  │  DockPulse Tab (React + MUI)         │  │
│  │  ┌─────────────┐  ┌───────────────┐  │  │
│  │  │  Container   │  │  Update       │  │  │
│  │  │  List w/     │  │  Details &    │  │  │
│  │  │  Status Dots │  │  History      │  │  │
│  │  └──────┬──────┘  └───────────────┘  │  │
│  │         │ socket                      │  │
│  │  ┌──────▼──────────────────────────┐  │  │
│  │  │  Backend Service (VM)           │  │  │
│  │  │  ┌────────────┐ ┌───────────┐  │  │  │
│  │  │  │ Scheduler  │ │ Registry  │  │  │  │
│  │  │  │ (cron)     │ │ Client    │  │  │  │
│  │  │  └────────────┘ └───────────┘  │  │  │
│  │  │  ┌────────────────────────────┐│  │  │
│  │  │  │ SQLite (check history,    ││  │  │
│  │  │  │ preferences, cache)       ││  │  │
│  │  │  └────────────────────────────┘│  │  │
│  │  └────────────────────────────────┘  │  │
│  └───────────────────────────────────────┘  │
│  Docker Engine API (read-only)              │
└─────────────────────────────────────────────┘
```

## Data Model

### Container Check Record

| Field | Type | Description |
|-------|------|-------------|
| id | INTEGER | Auto-increment PK |
| container_name | TEXT | Container name |
| image_ref | TEXT | Full image reference (e.g., `nginx:latest`) |
| local_digest | TEXT | SHA256 digest of local image |
| remote_digest | TEXT | SHA256 digest from registry |
| status | TEXT | `up-to-date`, `update-available`, `check-failed`, `unknown` |
| checked_at | DATETIME | Timestamp of last check |
| registry | TEXT | `dockerhub`, `ghcr`, `ecr`, `acr`, `custom` |

### User Preferences

| Field | Type | Description |
|-------|------|-------------|
| check_interval_minutes | INTEGER | How often to auto-check (default: 60) |
| notify_on_update | BOOLEAN | Show DD notification on new updates |
| ignored_images | TEXT | JSON array of images to skip |

## Privilege Requirements

**None beyond standard extension access.** DockPulse only needs:

- Docker Engine API (read-only) — list containers, inspect images
- Outbound HTTPS — query registry APIs
- Volume mount — persist SQLite database

No `SYS_ADMIN`, `NET_ADMIN`, or elevated capabilities required.

## MVP Scope (v0.1.0)

1. List running containers with image references
2. Check Docker Hub for newer digests (manual trigger)
3. Display status: ✅ up-to-date / ⚠️ update available / ❌ check failed
4. Store check history in SQLite

## Post-MVP Roadmap

- Scheduled automatic checks (v0.2.0)
- DD toast notifications (v0.3.0)
- Multi-registry auth (GHCR, private registries) (v0.4.0)
- Changelog preview — show what changed between versions (v0.5.0)
- One-click `docker pull` to update image (v0.6.0)

## Competitive Landscape

| Tool | Approach | Docker Desktop Integration |
|------|----------|---------------------------|
| Watchtower | Auto-updates containers | None — standalone container, development stalled |
| WUD (What's Up Docker) | Web dashboard for updates | None — separate web UI on port 3000 |
| Diun | Notification-only | None — CLI/notification tool |
| Tugtainer | Selective update approval | None — separate web UI |
| **DockPulse** | **Native DD extension** | **Built-in tab, zero extra containers** |
