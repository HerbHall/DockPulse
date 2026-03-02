package checker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/HerbHall/DockPulse/backend/internal/docker"
	"github.com/HerbHall/DockPulse/backend/internal/imageref"
	"github.com/HerbHall/DockPulse/backend/internal/registry"
	"github.com/HerbHall/DockPulse/backend/internal/store"
)

// Checker orchestrates image update checks by comparing local container
// image digests against remote registry digests.
type Checker struct {
	store    *store.Store
	registry registry.Registry
	docker   docker.Client
}

// New creates a Checker wired to the given store, registry, and Docker client.
func New(s *store.Store, reg registry.Registry, dc docker.Client) *Checker {
	return &Checker{
		store:    s,
		registry: reg,
		docker:   dc,
	}
}

// CheckAll enumerates running containers and checks each one for image
// updates. Results are persisted to the store and returned.
func (c *Checker) CheckAll(ctx context.Context) ([]store.ImageCheck, error) {
	containers, err := c.docker.ContainerList(ctx)
	if err != nil {
		return nil, fmt.Errorf("checker: list containers: %w", err)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	results := make([]store.ImageCheck, 0, len(containers))

	for i := range containers {
		check := c.checkContainer(ctx, &containers[i], now)
		if err = c.store.SaveCheck(ctx, &check); err != nil {
			log.Printf("checker: save check for %s: %v", containers[i].Name, err)
		}
		results = append(results, check)
	}

	return results, nil
}

func (c *Checker) checkContainer(ctx context.Context, ci *docker.ContainerInfo, timestamp string) store.ImageCheck {
	check := store.ImageCheck{
		ContainerName: ci.Name,
		ContainerID:   ci.ID,
		ImageRef:      ci.ImageRef,
		CheckedAt:     timestamp,
		Registry:      "dockerhub",
	}

	ref, err := imageref.Parse(ci.ImageRef)
	if err != nil {
		check.Status = store.StatusCheckFailed
		log.Printf("checker: parse image ref %q: %v", ci.ImageRef, err)
		return check
	}

	if ref.IsDigest {
		check.Status = store.StatusUpToDate
		return check
	}

	localDigest, err := c.docker.ImageDigest(ctx, ci.ImageID)
	if err != nil {
		check.Status = store.StatusCheckFailed
		log.Printf("checker: local digest for %s: %v", ci.Name, err)
		return check
	}
	check.LocalDigest = localDigest

	remoteDigest, err := c.registry.GetDigest(ctx, ref)
	if err != nil {
		check.Status = store.StatusCheckFailed
		log.Printf("checker: remote digest for %s: %v", ci.Name, err)
		return check
	}
	check.RemoteDigest = remoteDigest

	if localDigest == remoteDigest {
		check.Status = store.StatusUpToDate
	} else {
		check.Status = store.StatusUpdateAvailable
	}

	return check
}
