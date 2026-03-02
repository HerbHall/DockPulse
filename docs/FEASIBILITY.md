# DockPulse — Docker Desktop Extension Feasibility Assessment

## Summary

**Verdict**: Feasible and recommended as second project after RunNotes.

**Concept**: Native container update checker inside Docker Desktop. Queries registries for digest changes on running containers. No dependency on dying external tools (Watchtower is effectively dead — 8+ months no commits, community migrating to WUD/Diun/Tugtainer).

**Market Gap**: Docker Desktop has ZERO native update-checking functionality. All existing tools (WUD, Diun, Tugtainer, dockcheck) run as separate containers with separate web UIs. A DD extension providing update checks natively fills a legitimate gap.

## Why NOT Wrap Watchtower

- Watchtower development stalled, community fork (nicholas-fedor) is more active but still declining
- WUD (What's Up Docker) already has web UI on port 3000 — does what "Watchtower UI" would do
- Diun: notification-only, lightweight (<8MB RAM), no web UI
- Tugtainer: newest entry, web UI with selective update approval, multi-host agents
- dockcheck: simple bash script, zero overhead

Tying the extension to a dying tool is bad strategy. Build native capability instead.

## Competitive Landscape

| Tool | Web UI | In Docker Desktop | Update Check | Auto Update |
|------|--------|-------------------|--------------|-------------|
| Docker Desktop (native) | ✅ | ✅ | ❌ | ❌ |
| Watchtower | ❌ | ❌ | ✅ | ✅ |
| WUD | ✅ (port 3000) | ❌ | ✅ | Optional |
| Tugtainer | ✅ | ❌ | ✅ | Selective |
| Diun | ❌ | ❌ | ✅ | ❌ |
| **DockPulse** | **✅ (native tab)** | **✅** | **✅** | **Optional** |

## Technical Approach

- Extension backend queries Docker Hub / GHCR / registry APIs for digest comparison
- Compare running image digests against registry latest
- React frontend shows update availability per container in DD-native UI
- No external dependencies, works out of the box
- Same SDK pattern as RunNotes (React + backend service + socket communication)

## Effort Estimate

| Tier | Scope | Effort |
|------|-------|--------|
| MVP | Query registries for digest changes, display status per container | 3-5 days |
| Solid | Notifications, update history, schedule controls | 1 week |
| Full | Multi-registry auth, changelog links, rollback capability | 2-3 weeks |

## Risks

1. **Registry API complexity** — Docker Hub, GHCR, private registries have varying auth and rate limits
2. **Update safety** — Auto-updating is inherently dangerous; clear safeguards needed
3. **Rate limiting** — Docker Hub anonymous pull limit (100/6hrs) could be hit with many containers

## Inspiration

- WUD dashboard: container list with version columns, update-kind indicators (major/minor/patch)
- Tugtainer: per-container approve/reject, dependency-aware restart ordering
- Renovate/Dependabot model: "show available updates, let the human decide"
- UnRAID Community Apps: shows available updates with changelogs inline
