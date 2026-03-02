# Contributing to DockPulse

Thanks for your interest in contributing!

## Getting Started

1. Fork the repository
2. Create a feature branch (`git checkout -b feat/my-feature`)
3. Make your changes
4. Commit using conventional commits (`feat:`, `fix:`, `docs:`, `chore:`)
5. Push and open a Pull Request

## Development Setup

Requires Docker Desktop 4.8.0+ with Extensions enabled.

```bash
make build-extension    # Build
make install-extension  # Install into Docker Desktop
make update-extension   # Update after changes
```

## Code Style

- See `.editorconfig` for formatting rules
- Frontend: React + Material UI
- Conventional commits required
