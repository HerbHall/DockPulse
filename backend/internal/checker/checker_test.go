package checker

import (
	"context"
	"errors"
	"testing"

	"github.com/HerbHall/DockPulse/backend/internal/docker"
	"github.com/HerbHall/DockPulse/backend/internal/imageref"
	"github.com/HerbHall/DockPulse/backend/internal/store"
)

// mockRegistry implements registry.Registry for testing.
type mockRegistry struct {
	digests map[string]string // "namespace/name:tag" -> digest
	err     error
}

func (m *mockRegistry) GetDigest(_ context.Context, ref *imageref.ImageRef) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	key := ref.FullRepository() + ":" + ref.Tag
	digest, ok := m.digests[key]
	if !ok {
		return "", errors.New("mock: digest not found")
	}
	return digest, nil
}

// mockDockerClient implements docker.Client for testing.
type mockDockerClient struct {
	containers []docker.ContainerInfo
	digests    map[string]string // imageID -> local digest
	listErr    error
	digestErr  error
}

func (m *mockDockerClient) ContainerList(_ context.Context) ([]docker.ContainerInfo, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.containers, nil
}

func (m *mockDockerClient) ImageDigest(_ context.Context, imageID string) (string, error) {
	if m.digestErr != nil {
		return "", m.digestErr
	}
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

func TestCheckAll_UpdateAvailable(t *testing.T) {
	s := newTestStore(t)
	dc := &mockDockerClient{
		containers: []docker.ContainerInfo{
			{ID: "c1", Name: "web", ImageRef: "nginx:latest", ImageID: "img1", State: "running"},
		},
		digests: map[string]string{"img1": "sha256:old"},
	}
	reg := &mockRegistry{
		digests: map[string]string{"library/nginx:latest": "sha256:new"},
	}

	chk := New(s, reg, dc)
	results, err := chk.CheckAll(context.Background())
	if err != nil {
		t.Fatalf("CheckAll: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != store.StatusUpdateAvailable {
		t.Errorf("status = %q, want %q", results[0].Status, store.StatusUpdateAvailable)
	}
	if results[0].LocalDigest != "sha256:old" {
		t.Errorf("local digest = %q, want %q", results[0].LocalDigest, "sha256:old")
	}
	if results[0].RemoteDigest != "sha256:new" {
		t.Errorf("remote digest = %q, want %q", results[0].RemoteDigest, "sha256:new")
	}
}

func TestCheckAll_UpToDate(t *testing.T) {
	s := newTestStore(t)
	dc := &mockDockerClient{
		containers: []docker.ContainerInfo{
			{ID: "c1", Name: "web", ImageRef: "nginx:latest", ImageID: "img1", State: "running"},
		},
		digests: map[string]string{"img1": "sha256:same"},
	}
	reg := &mockRegistry{
		digests: map[string]string{"library/nginx:latest": "sha256:same"},
	}

	chk := New(s, reg, dc)
	results, err := chk.CheckAll(context.Background())
	if err != nil {
		t.Fatalf("CheckAll: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != store.StatusUpToDate {
		t.Errorf("status = %q, want %q", results[0].Status, store.StatusUpToDate)
	}
}

func TestCheckAll_DigestPinned(t *testing.T) {
	s := newTestStore(t)
	dc := &mockDockerClient{
		containers: []docker.ContainerInfo{
			{ID: "c1", Name: "web", ImageRef: "nginx@sha256:abc123", ImageID: "img1", State: "running"},
		},
		digests: map[string]string{"img1": "sha256:abc123"},
	}
	reg := &mockRegistry{
		digests: map[string]string{},
	}

	chk := New(s, reg, dc)
	results, err := chk.CheckAll(context.Background())
	if err != nil {
		t.Fatalf("CheckAll: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != store.StatusUpToDate {
		t.Errorf("status = %q, want %q", results[0].Status, store.StatusUpToDate)
	}
}

func TestCheckAll_RegistryError(t *testing.T) {
	s := newTestStore(t)
	dc := &mockDockerClient{
		containers: []docker.ContainerInfo{
			{ID: "c1", Name: "web", ImageRef: "nginx:latest", ImageID: "img1", State: "running"},
		},
		digests: map[string]string{"img1": "sha256:local"},
	}
	reg := &mockRegistry{
		err: errors.New("registry unavailable"),
	}

	chk := New(s, reg, dc)
	results, err := chk.CheckAll(context.Background())
	if err != nil {
		t.Fatalf("CheckAll: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != store.StatusCheckFailed {
		t.Errorf("status = %q, want %q", results[0].Status, store.StatusCheckFailed)
	}
}

func TestCheckAll_MultipleContainers(t *testing.T) {
	s := newTestStore(t)
	dc := &mockDockerClient{
		containers: []docker.ContainerInfo{
			{ID: "c1", Name: "web", ImageRef: "nginx:latest", ImageID: "img1", State: "running"},
			{ID: "c2", Name: "api", ImageRef: "myuser/app:v2", ImageID: "img2", State: "running"},
			{ID: "c3", Name: "db", ImageRef: "postgres@sha256:pinned", ImageID: "img3", State: "running"},
		},
		digests: map[string]string{
			"img1": "sha256:old",
			"img2": "sha256:current",
			"img3": "sha256:pinned",
		},
	}
	reg := &mockRegistry{
		digests: map[string]string{
			"library/nginx:latest": "sha256:new",
			"myuser/app:v2":        "sha256:current",
		},
	}

	chk := New(s, reg, dc)
	results, err := chk.CheckAll(context.Background())
	if err != nil {
		t.Fatalf("CheckAll: %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	expectations := []struct {
		name   string
		status store.CheckStatus
	}{
		{"web", store.StatusUpdateAvailable},
		{"api", store.StatusUpToDate},
		{"db", store.StatusUpToDate},
	}

	for i, exp := range expectations {
		if results[i].ContainerName != exp.name {
			t.Errorf("result[%d] name = %q, want %q", i, results[i].ContainerName, exp.name)
		}
		if results[i].Status != exp.status {
			t.Errorf("result[%d] (%s) status = %q, want %q", i, exp.name, results[i].Status, exp.status)
		}
	}

	// Verify results were persisted to the store.
	stored, err := s.GetLatestChecks(context.Background())
	if err != nil {
		t.Fatalf("GetLatestChecks: %v", err)
	}
	if len(stored) != 3 {
		t.Errorf("stored checks = %d, want 3", len(stored))
	}
}
