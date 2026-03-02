# Changelog

All notable changes to DockPulse will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Go backend with Docker Hub v2 registry client (anonymous token auth, manifest HEAD for digest comparison)
- Image reference parser supporting Docker Hub, GHCR, and private registry formats
- SQLite data model for storing image check results (`modernc.org/sqlite`, pure Go)
- Checker orchestrator: enumerates running containers, compares digests, stores results
- HTTP API: `GET /api/checks`, `POST /api/check-all`, `GET /api/status`
- React + MUI v5 frontend with container table, color-coded status chips, and error boundary
- `useBackend` hook wiring frontend to backend via Docker Desktop extension SDK
- Multi-stage Dockerfile (Go 1.24-alpine + Node 22-alpine + Alpine 3.19)
- GitHub Actions CI with 7 parallel jobs (go-lint, go-build, go-test, fe-lint, fe-typecheck, fe-test, docker-build)
- Makefile targets: `validate`, `go-test`, `go-lint`, `fe-test`, `fe-typecheck`, `fe-lint`
- Initial project scaffold
