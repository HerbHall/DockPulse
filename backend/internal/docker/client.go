package docker

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types/container"
	dockerclient "github.com/docker/docker/client"
)

// ContainerInfo holds the essential details of a running container.
type ContainerInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ImageRef string `json:"imageRef"`
	ImageID  string `json:"imageId"`
	State    string `json:"state"`
}

// Client defines the interface for Docker Engine interactions,
// allowing mock implementations in tests.
type Client interface {
	ContainerList(ctx context.Context) ([]ContainerInfo, error)
	ImageDigest(ctx context.Context, imageID string) (string, error)
	Close() error
}

// DockerClient wraps the Docker Engine SDK client.
type DockerClient struct {
	cli *dockerclient.Client
}

// NewClient creates a Docker client using environment configuration
// and automatic API version negotiation.
func NewClient() (*DockerClient, error) {
	cli, err := dockerclient.NewClientWithOpts(
		dockerclient.FromEnv,
		dockerclient.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, fmt.Errorf("docker: create client: %w", err)
	}
	return &DockerClient{cli: cli}, nil
}

// ContainerList returns information about all running containers.
func (d *DockerClient) ContainerList(ctx context.Context) ([]ContainerInfo, error) {
	containers, err := d.cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("docker: list containers: %w", err)
	}

	result := make([]ContainerInfo, 0, len(containers))
	for i := range containers {
		name := ""
		if len(containers[i].Names) > 0 {
			name = strings.TrimPrefix(containers[i].Names[0], "/")
		}

		result = append(result, ContainerInfo{
			ID:       containers[i].ID,
			Name:     name,
			ImageRef: containers[i].Image,
			ImageID:  containers[i].ImageID,
			State:    containers[i].State,
		})
	}

	return result, nil
}

// ImageDigest inspects the given image and returns the first RepoDigest
// if available. Returns an empty string when no digest is found.
func (d *DockerClient) ImageDigest(ctx context.Context, imageID string) (string, error) {
	inspect, err := d.cli.ImageInspect(ctx, imageID)
	if err != nil {
		return "", fmt.Errorf("docker: inspect image %s: %w", imageID, err)
	}

	if len(inspect.RepoDigests) > 0 {
		// RepoDigests are in the form "repo@sha256:abc...", extract the digest.
		parts := strings.SplitN(inspect.RepoDigests[0], "@", 2)
		if len(parts) == 2 {
			return parts[1], nil
		}
		return inspect.RepoDigests[0], nil
	}

	return "", nil
}

// Close releases the Docker client resources.
func (d *DockerClient) Close() error {
	return d.cli.Close()
}

// Compile-time interface check.
var _ Client = (*DockerClient)(nil)
