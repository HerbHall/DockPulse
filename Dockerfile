# syntax=docker/dockerfile:1

# DockPulse Docker Desktop Extension
# Multi-stage build: backend + frontend

FROM alpine
LABEL org.opencontainers.image.title="DockPulse" \
      org.opencontainers.image.description="Check if your container images have newer versions available" \
      org.opencontainers.image.vendor="Herb Hall" \
      com.docker.desktop.extension.api.version=">= 0.3.3" \
      com.docker.desktop.extension.icon="https://raw.githubusercontent.com/HerbHall/DockPulse/main/docker.svg" \
      com.docker.extension.screenshots="" \
      com.docker.extension.detailed-description="DockPulse checks your running containers against their registries and shows which images have updates available — right inside Docker Desktop." \
      com.docker.extension.publisher-url="https://github.com/HerbHall" \
      com.docker.extension.changelog=""

# TODO: Add backend build stage
# TODO: Add frontend build stage
# TODO: Copy metadata.json and ui assets

COPY metadata.json .
