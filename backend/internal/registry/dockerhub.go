package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/HerbHall/DockPulse/backend/internal/imageref"
)

// DockerHubClient implements the Registry interface by querying the
// Docker Hub v2 API for image manifest digests.
type DockerHubClient struct {
	httpClient  *http.Client
	authURL     string
	registryURL string
}

// NewDockerHubClient returns a client configured for the public Docker Hub.
func NewDockerHubClient() *DockerHubClient {
	return &DockerHubClient{
		httpClient:  &http.Client{Timeout: 30 * time.Second},
		authURL:     "https://auth.docker.io",
		registryURL: "https://registry-1.docker.io",
	}
}

// tokenResponse holds the JSON body returned by the Docker auth endpoint.
type tokenResponse struct {
	Token string `json:"token"`
}

// getToken requests a pull-scoped bearer token for the given repository.
func (c *DockerHubClient) getToken(ctx context.Context, repository string) (string, error) {
	url := fmt.Sprintf("%s/token?service=registry.docker.io&scope=repository:%s:pull", c.authURL, repository)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return "", fmt.Errorf("dockerhub: build token request: %w", err)
	}

	resp, err := c.httpClient.Do(req) //nolint:gosec // G107: URL built from trusted authURL config
	if err != nil {
		return "", fmt.Errorf("dockerhub: failed to get token for %s: %w", repository, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("dockerhub: auth returned status %d for %s", resp.StatusCode, repository)
	}

	var tr tokenResponse
	if err = json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return "", fmt.Errorf("dockerhub: decode token response for %s: %w", repository, err)
	}

	if tr.Token == "" {
		return "", fmt.Errorf("dockerhub: empty token for %s", repository)
	}

	return tr.Token, nil
}

// GetDigest fetches the remote manifest digest for the given image reference
// using a HEAD request to avoid counting against Docker Hub pull rate limits.
func (c *DockerHubClient) GetDigest(ctx context.Context, ref *imageref.ImageRef) (string, error) {
	repo := ref.FullRepository()
	tag := ref.Tag
	if tag == "" {
		tag = "latest"
	}

	token, err := c.getToken(ctx, repo)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/v2/%s/manifests/%s", c.registryURL, repo, tag)

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, http.NoBody)
	if err != nil {
		return "", fmt.Errorf("dockerhub: build manifest request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json, application/vnd.oci.image.index.v1+json")

	resp, err := c.httpClient.Do(req) //nolint:gosec // G107: URL built from trusted registryURL config
	if err != nil {
		return "", fmt.Errorf("dockerhub: manifest request failed for %s:%s: %w", repo, tag, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("dockerhub: manifest not found for %s:%s", repo, tag)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("dockerhub: manifest request returned status %d for %s:%s", resp.StatusCode, repo, tag)
	}

	digest := resp.Header.Get("Docker-Content-Digest")
	if digest == "" {
		return "", fmt.Errorf("dockerhub: no Docker-Content-Digest header for %s:%s", repo, tag)
	}

	return digest, nil
}
