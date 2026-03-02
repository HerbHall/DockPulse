# DockPulse

Docker Desktop extension that checks if your running container images have newer versions available.

## Why

Docker Desktop shows container names, images, status, and ports — but never tells you if a newer image version exists. DockPulse adds a native update-checking tab so you always know which containers are running stale images.

## Features (Planned)

- Container list with update status badges (up-to-date / update available / unknown)
- Registry digest comparison against Docker Hub, GHCR, and private registries
- Scheduled background checks with configurable intervals
- Update history log
- One-click pull + recreate (with safety confirmation)

## Status

Pre-development scaffold. See [HANDOFF.md](HANDOFF.md) for next steps.

## Development

```bash
make build-extension    # Build the extension image
make install-extension  # Install into Docker Desktop
make update-extension   # Update after changes
```

## License

[MIT](LICENSE) © 2026 Herb Hall
