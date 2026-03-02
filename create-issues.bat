@echo off
cd /d D:\devspace\DockPulse

echo Creating issue 1...
echo Design the data model for tracking container image update status.> .body
echo.>> .body
echo Fields: container_name, container_id, image_reference, local_digest, registry_digest, update_available, first_seen, last_checked, check_count.>> .body
echo.>> .body
echo Storage on Docker volume (SQLite or JSON). See docs/ARCHITECTURE.md.>> .body
gh issue create --title "feat: data model - image update tracking schema" --body-file .body --label "feat,mvp"

echo Creating issue 2...
echo Build the React frontend tab for Docker Desktop.> .body
echo.>> .body
echo Components: container list with update status indicators (current/update-available/unknown/error), detail panel with digest info, manual Check Now button.>> .body
echo.>> .body
echo Must use Material UI to match Docker Desktop look. All state fetched from backend on mount.>> .body
gh issue create --title "feat: React UI - update status dashboard" --body-file .body --label "feat,mvp"

echo Creating issue 3...
echo Build the backend service that queries registries for image updates.> .body
echo.>> .body
echo Responsibilities: enumerate running containers, resolve image tags to registry digests via Registry HTTP API v2, compare local vs registry digests, store results.>> .body
echo.>> .body
echo Communication via Extensions SDK socket. Docker Hub as Tier 1 registry.>> .body
gh issue create --title "feat: backend service - registry digest checker" --body-file .body --label "feat,mvp"

echo Creating issue 4...
echo Parse image references into registry, repository, and tag components.> .body
echo.>> .body
echo Handle: official images (postgres:16), namespaced images (library/postgres), custom registries (ghcr.io/user/repo:tag), digest-pinned images (repo@sha256:...), images with no explicit tag (defaults to latest).>> .body
gh issue create --title "feat: image reference parsing and normalization" --body-file .body --label "feat,mvp"

echo Creating issue 5...
echo Implement Docker Hub authentication and rate limit management.> .body
echo.>> .body
echo Docker Hub allows 100 pulls/6h anonymous, 200 authenticated. Use manifest HEAD requests where possible. Cache results to minimize API calls. Support Docker Hub auth tokens.>> .body
gh issue create --title "feat: Docker Hub auth and rate limit handling" --body-file .body --label "feat,mvp"

echo Creating issue 6...
echo Add scheduled background scanning on a configurable interval.> .body
echo.>> .body
echo Default: check every 6 hours. User-configurable via settings panel. Include last-scan timestamp in UI. Manual Check Now always available.>> .body
gh issue create --title "feat: scheduled scan interval with settings" --body-file .body --label "feat,enhancement"

echo Creating issue 7...
echo Store and display scan history showing when updates were detected.> .body
echo.>> .body
echo Track per-container: when update first appeared, how many checks since, whether user acknowledged it. Provide a history view in the UI.>> .body
gh issue create --title "feat: scan history and update timeline" --body-file .body --label "feat,enhancement"

echo Creating issue 8...
echo Add support for GitHub Container Registry (GHCR) and other OCI registries.> .body
echo.>> .body
echo Extend the registry checker beyond Docker Hub. GHCR uses token-based auth with ghcr.io. Follow OCI Distribution Spec for generic registry support.>> .body
gh issue create --title "feat: multi-registry support (GHCR, Quay, private)" --body-file .body --label "feat,enhancement"

echo Creating issue 9...
echo Publish extension to Docker Hub and submit to the Extensions Marketplace.> .body
echo.>> .body
echo Requires: multi-arch build (linux/amd64 + linux/arm64), proper OCI labels, icon, screenshots, detailed description.>> .body
echo.>> .body
echo Image name: herbhall/dockpulse>> .body
gh issue create --title "chore: Docker Hub publishing and marketplace submission" --body-file .body --label "chore"

echo Creating issue 10...
echo Design the extension icon and branding assets.> .body
echo.>> .body
echo Need: SVG icon for Docker Desktop sidebar, screenshot images for marketplace listing, color scheme (green/earthy preferred). Name: DockPulse - "check the pulse of your containers".>> .body
gh issue create --title "docs: extension icon and branding" --body-file .body --label "docs"

del .body
echo Done - all issues created.
