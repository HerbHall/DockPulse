# syntax=docker/dockerfile:1

# DockPulse Docker Desktop Extension
# Multi-stage build: Go backend + React frontend

ARG VERSION=0.1.0 # x-release-please-version

# ---------- Stage 1: Build Go backend ----------
FROM golang:1.24-alpine AS backend-build
WORKDIR /build

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o dockpulse .

# ---------- Stage 2: Build React frontend ----------
FROM node:22-alpine AS frontend-build
WORKDIR /build

COPY ui/package.json ui/package-lock.json ./
RUN npm ci

COPY ui/ .
RUN npm run build

# ---------- Stage 3: Final image ----------
FROM alpine:3.19

LABEL org.opencontainers.image.title="DockPulse" \
      org.opencontainers.image.description="Check if your container images have newer versions available" \
      org.opencontainers.image.vendor="Herb Hall" \
      com.docker.desktop.extension.api.version=">= 0.3.3" \
      com.docker.desktop.extension.icon="docker.svg" \
      com.docker.extension.detailed-description="DockPulse checks your running containers against their registries and shows which images have updates available — right inside Docker Desktop." \
      com.docker.extension.publisher-url="https://github.com/HerbHall" \
      com.docker.extension.screenshots="" \
      com.docker.extension.changelog=""

COPY metadata.json /
COPY docker.svg /
COPY --from=backend-build /build/dockpulse /
COPY --from=frontend-build /build/build /ui

ENTRYPOINT ["/dockpulse"]
