package registry

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HerbHall/DockPulse/backend/internal/imageref"
)

func newTestClient(authURL, registryURL string) *DockerHubClient {
	return &DockerHubClient{
		httpClient:  http.DefaultClient,
		authURL:     authURL,
		registryURL: registryURL,
	}
}

func TestGetDigest_Success(t *testing.T) {
	wantDigest := "sha256:abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"

	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{"token":"test-token-abc"}`)
	}))
	defer authServer.Close()

	registryServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodHead {
			t.Errorf("expected HEAD request, got %s", r.Method)
		}
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token-abc" {
			t.Errorf("expected Bearer token, got %q", auth)
		}
		w.Header().Set("Docker-Content-Digest", wantDigest)
		w.WriteHeader(http.StatusOK)
	}))
	defer registryServer.Close()

	client := newTestClient(authServer.URL, registryServer.URL)
	ref := &imageref.ImageRef{
		Registry:  "registry-1.docker.io",
		Namespace: "library",
		Name:      "nginx",
		Tag:       "latest",
	}

	got, err := client.GetDigest(context.Background(), ref)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != wantDigest {
		t.Errorf("digest = %q, want %q", got, wantDigest)
	}
}

func TestGetDigest_AuthFailure(t *testing.T) {
	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer authServer.Close()

	registryServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		t.Fatal("registry should not be called when auth fails")
	}))
	defer registryServer.Close()

	client := newTestClient(authServer.URL, registryServer.URL)
	ref := &imageref.ImageRef{
		Registry:  "registry-1.docker.io",
		Namespace: "library",
		Name:      "nginx",
		Tag:       "latest",
	}

	_, err := client.GetDigest(context.Background(), ref)
	if err == nil {
		t.Fatal("expected error for auth failure, got nil")
	}
}

func TestGetDigest_ManifestNotFound(t *testing.T) {
	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{"token":"test-token"}`)
	}))
	defer authServer.Close()

	registryServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer registryServer.Close()

	client := newTestClient(authServer.URL, registryServer.URL)
	ref := &imageref.ImageRef{
		Registry:  "registry-1.docker.io",
		Namespace: "library",
		Name:      "nonexistent",
		Tag:       "latest",
	}

	_, err := client.GetDigest(context.Background(), ref)
	if err == nil {
		t.Fatal("expected error for manifest not found, got nil")
	}
}

func TestGetDigest_NoDigestHeader(t *testing.T) {
	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{"token":"test-token"}`)
	}))
	defer authServer.Close()

	registryServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer registryServer.Close()

	client := newTestClient(authServer.URL, registryServer.URL)
	ref := &imageref.ImageRef{
		Registry:  "registry-1.docker.io",
		Namespace: "library",
		Name:      "nginx",
		Tag:       "latest",
	}

	_, err := client.GetDigest(context.Background(), ref)
	if err == nil {
		t.Fatal("expected error for missing digest header, got nil")
	}
}
