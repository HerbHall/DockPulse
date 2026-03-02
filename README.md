# DockPulse

Docker Desktop extension that checks if your running container images have newer versions available.

## Why

Docker Desktop shows container names, images, status, and ports — but never tells you if a newer image version exists. DockPulse adds a native update-checking tab so you always know which containers are running stale images.

## Features

- Container list with update status badges (up-to-date / update available / check failed / unknown)
- Registry digest comparison against Docker Hub
- One-click "Check Now" to scan all running containers
- Color-coded status chips for quick visual triage
- Error handling with dismissible alerts

### Planned

- GHCR and private registry support ([#8](https://github.com/HerbHall/DockPulse/issues/8))
- Scheduled background checks with configurable intervals ([#6](https://github.com/HerbHall/DockPulse/issues/6))
- Scan history and update timeline ([#7](https://github.com/HerbHall/DockPulse/issues/7))
- Docker Hub marketplace publishing ([#9](https://github.com/HerbHall/DockPulse/issues/9))

## How It Works

1. Lists all running containers via the Docker Engine API
2. Parses each container's image reference (registry, namespace, name, tag)
3. Fetches the latest manifest digest from Docker Hub using HEAD requests (avoids counting against pull rate limits)
4. Compares the remote digest against the local image's `RepoDigests`
5. Stores results in SQLite and displays status in the extension UI

## Development

```bash
make validate            # Run all checks (build, test, lint, typecheck)
make build-extension     # Build the extension image
make install-extension   # Install into Docker Desktop
make update-extension    # Update after code changes

# Individual targets
make go-test             # Backend tests
make go-lint             # Backend lint (golangci-lint v2)
make fe-test             # Frontend tests (vitest)
make fe-typecheck        # TypeScript type checking
make fe-lint             # Frontend lint (eslint)
```

## Tech Stack

- **Frontend**: React 18 + MUI v5 + Vite 7 + TypeScript
- **Backend**: Go 1.24 + stdlib `net/http` + Docker Engine API
- **Storage**: SQLite via `modernc.org/sqlite` (pure Go, no CGO)
- **Build**: Multi-stage Dockerfile (Go + Node + Alpine)

## License

[MIT](LICENSE) © 2026 Herb Hall
