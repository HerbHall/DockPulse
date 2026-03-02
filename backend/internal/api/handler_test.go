package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HerbHall/DockPulse/backend/internal/checker"
	"github.com/HerbHall/DockPulse/backend/internal/docker"
	"github.com/HerbHall/DockPulse/backend/internal/imageref"
	"github.com/HerbHall/DockPulse/backend/internal/store"
)

// mockRegistry implements registry.Registry for handler tests.
type mockRegistry struct {
	digests map[string]string
}

func (m *mockRegistry) GetDigest(_ context.Context, ref *imageref.ImageRef) (string, error) {
	key := ref.FullRepository() + ":" + ref.Tag
	return m.digests[key], nil
}

// mockDockerClient implements docker.Client for handler tests.
type mockDockerClient struct {
	containers []docker.ContainerInfo
	digests    map[string]string
}

func (m *mockDockerClient) ContainerList(_ context.Context) ([]docker.ContainerInfo, error) {
	return m.containers, nil
}

func (m *mockDockerClient) ImageDigest(_ context.Context, imageID string) (string, error) {
	return m.digests[imageID], nil
}

func (m *mockDockerClient) Close() error {
	return nil
}

func newTestStore(t *testing.T) *store.Store {
	t.Helper()
	s, err := store.New(":memory:")
	if err != nil {
		t.Fatalf("new store: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

func setupHandler(t *testing.T) (*Handler, *http.ServeMux) {
	t.Helper()
	s := newTestStore(t)
	dc := &mockDockerClient{
		containers: []docker.ContainerInfo{
			{ID: "c1", Name: "web", ImageRef: "nginx:latest", ImageID: "img1", State: "running"},
		},
		digests: map[string]string{"img1": "sha256:local"},
	}
	reg := &mockRegistry{
		digests: map[string]string{"library/nginx:latest": "sha256:remote"},
	}
	chk := checker.New(s, reg, dc)
	h := NewHandler(chk, s)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return h, mux
}

func TestGetChecks_Empty(t *testing.T) {
	_, mux := setupHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/checks", http.NoBody)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("Content-Type = %q, want %q", ct, "application/json")
	}

	var resp checksResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(resp.Checks) != 0 {
		t.Errorf("expected 0 checks, got %d", len(resp.Checks))
	}
}

func TestCheckAll(t *testing.T) {
	_, mux := setupHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/api/check-all", http.NoBody)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp checkAllResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if resp.StartedAt == "" {
		t.Error("expected non-empty startedAt")
	}
	if len(resp.Checks) != 1 {
		t.Fatalf("expected 1 check, got %d", len(resp.Checks))
	}
	if resp.Checks[0].Status != store.StatusUpdateAvailable {
		t.Errorf("status = %q, want %q", resp.Checks[0].Status, store.StatusUpdateAvailable)
	}
}

func TestGetChecks_AfterCheckAll(t *testing.T) {
	_, mux := setupHandler(t)

	// Trigger a check first.
	req := httptest.NewRequest(http.MethodPost, "/api/check-all", http.NoBody)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("check-all status = %d, want %d", rec.Code, http.StatusOK)
	}

	// Now GET /api/checks should return the persisted results.
	req = httptest.NewRequest(http.MethodGet, "/api/checks", http.NoBody)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("get checks status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp checksResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(resp.Checks) != 1 {
		t.Fatalf("expected 1 check, got %d", len(resp.Checks))
	}
	if resp.Checks[0].ContainerName != "web" {
		t.Errorf("container name = %q, want %q", resp.Checks[0].ContainerName, "web")
	}
}

func TestStatus(t *testing.T) {
	_, mux := setupHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/status", http.NoBody)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp statusResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !resp.Healthy {
		t.Error("expected healthy = true")
	}
	if resp.Version != "0.1.0" {
		t.Errorf("version = %q, want %q", resp.Version, "0.1.0")
	}
}
