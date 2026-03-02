package store_test

import (
	"context"
	"testing"

	"github.com/HerbHall/DockPulse/backend/internal/store"
)

func newTestStore(t *testing.T) *store.Store {
	t.Helper()
	s, err := store.New(":memory:")
	if err != nil {
		t.Fatalf("new store: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

func TestSaveAndGetCheck(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	check := &store.ImageCheck{
		ContainerName: "web-app",
		ContainerID:   "abc123",
		ImageRef:      "nginx:latest",
		LocalDigest:   "sha256:aaa",
		RemoteDigest:  "sha256:bbb",
		Status:        store.StatusUpdateAvailable,
		CheckedAt:     "2026-03-02T10:00:00Z",
		Registry:      "dockerhub",
	}

	if err := s.SaveCheck(ctx, check); err != nil {
		t.Fatalf("save check: %v", err)
	}

	if check.ID == 0 {
		t.Fatal("expected non-zero ID after save")
	}

	got, err := s.GetCheckByContainerID(ctx, "abc123")
	if err != nil {
		t.Fatalf("get check: %v", err)
	}
	if got == nil {
		t.Fatal("expected non-nil check")
	}

	tests := []struct {
		name string
		got  string
		want string
	}{
		{"ContainerName", got.ContainerName, "web-app"},
		{"ContainerID", got.ContainerID, "abc123"},
		{"ImageRef", got.ImageRef, "nginx:latest"},
		{"LocalDigest", got.LocalDigest, "sha256:aaa"},
		{"RemoteDigest", got.RemoteDigest, "sha256:bbb"},
		{"Status", string(got.Status), string(store.StatusUpdateAvailable)},
		{"CheckedAt", got.CheckedAt, "2026-03-02T10:00:00Z"},
		{"Registry", got.Registry, "dockerhub"},
	}
	for _, tc := range tests {
		if tc.got != tc.want {
			t.Errorf("%s: got %q, want %q", tc.name, tc.got, tc.want)
		}
	}
}

func TestGetLatestChecks(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	checks := []*store.ImageCheck{
		{
			ContainerName: "web", ContainerID: "c1", ImageRef: "nginx:1.24",
			Status: store.StatusUpToDate, CheckedAt: "2026-03-01T10:00:00Z", Registry: "dockerhub",
		},
		{
			ContainerName: "web", ContainerID: "c1", ImageRef: "nginx:1.24",
			Status: store.StatusUpdateAvailable, CheckedAt: "2026-03-02T10:00:00Z", Registry: "dockerhub",
		},
		{
			ContainerName: "api", ContainerID: "c2", ImageRef: "myapp:latest",
			Status: store.StatusUpToDate, CheckedAt: "2026-03-02T09:00:00Z", Registry: "dockerhub",
		},
	}
	for _, c := range checks {
		if err := s.SaveCheck(ctx, c); err != nil {
			t.Fatalf("save check: %v", err)
		}
	}

	latest, err := s.GetLatestChecks(ctx)
	if err != nil {
		t.Fatalf("get latest checks: %v", err)
	}

	if got := len(latest); got != 2 {
		t.Fatalf("expected 2 latest checks (one per container), got %d", got)
	}

	// Results ordered by checked_at DESC: c1 (03-02) before c2 (03-02T09)
	if latest[0].ContainerID != "c1" {
		t.Errorf("first result: got container %q, want %q", latest[0].ContainerID, "c1")
	}
	if latest[0].Status != store.StatusUpdateAvailable {
		t.Errorf("first result status: got %q, want %q", latest[0].Status, store.StatusUpdateAvailable)
	}
	if latest[1].ContainerID != "c2" {
		t.Errorf("second result: got container %q, want %q", latest[1].ContainerID, "c2")
	}
}

func TestGetCheckByContainerID(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	tests := []struct {
		name        string
		containerID string
		wantNil     bool
	}{
		{"existing container", "existing-id", false},
		{"non-existing container", "no-such-id", true},
	}

	// Seed one check.
	if err := s.SaveCheck(ctx, &store.ImageCheck{
		ContainerName: "db", ContainerID: "existing-id", ImageRef: "postgres:16",
		Status: store.StatusUpToDate, CheckedAt: "2026-03-02T08:00:00Z", Registry: "dockerhub",
	}); err != nil {
		t.Fatalf("save check: %v", err)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := s.GetCheckByContainerID(ctx, tc.containerID)
			if err != nil {
				t.Fatalf("get check: %v", err)
			}
			if tc.wantNil && got != nil {
				t.Fatalf("expected nil, got %+v", got)
			}
			if !tc.wantNil && got == nil {
				t.Fatal("expected non-nil check")
			}
		})
	}
}

func TestPreferences(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	// Missing key returns empty string.
	val, err := s.GetPreference(ctx, "missing-key")
	if err != nil {
		t.Fatalf("get missing preference: %v", err)
	}
	if val != "" {
		t.Errorf("missing key: got %q, want empty string", val)
	}

	// Set and retrieve.
	if err := s.SetPreference(ctx, "check_interval", "30m"); err != nil {
		t.Fatalf("set preference: %v", err)
	}
	val, err = s.GetPreference(ctx, "check_interval")
	if err != nil {
		t.Fatalf("get preference: %v", err)
	}
	if val != "30m" {
		t.Errorf("check_interval: got %q, want %q", val, "30m")
	}

	// Overwrite existing.
	if err := s.SetPreference(ctx, "check_interval", "1h"); err != nil {
		t.Fatalf("set preference overwrite: %v", err)
	}
	val, err = s.GetPreference(ctx, "check_interval")
	if err != nil {
		t.Fatalf("get preference after overwrite: %v", err)
	}
	if val != "1h" {
		t.Errorf("check_interval after overwrite: got %q, want %q", val, "1h")
	}
}

func TestMultipleChecksForSameContainer(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	// Save three checks for the same container at different times.
	timestamps := []string{
		"2026-03-02T08:00:00Z",
		"2026-03-02T10:00:00Z",
		"2026-03-02T09:00:00Z",
	}
	statuses := []store.CheckStatus{
		store.StatusUpToDate,
		store.StatusUpdateAvailable,
		store.StatusCheckFailed,
	}

	for i, ts := range timestamps {
		if err := s.SaveCheck(ctx, &store.ImageCheck{
			ContainerName: "app", ContainerID: "repeat-id", ImageRef: "myapp:latest",
			Status: statuses[i], CheckedAt: ts, Registry: "dockerhub",
		}); err != nil {
			t.Fatalf("save check %d: %v", i, err)
		}
	}

	// GetCheckByContainerID returns the most recent by checked_at.
	got, err := s.GetCheckByContainerID(ctx, "repeat-id")
	if err != nil {
		t.Fatalf("get check: %v", err)
	}
	if got == nil {
		t.Fatal("expected non-nil check")
	}
	if got.CheckedAt != "2026-03-02T10:00:00Z" {
		t.Errorf("CheckedAt: got %q, want %q", got.CheckedAt, "2026-03-02T10:00:00Z")
	}
	if got.Status != store.StatusUpdateAvailable {
		t.Errorf("Status: got %q, want %q", got.Status, store.StatusUpdateAvailable)
	}

	// GetLatestChecks should return only one entry for this container.
	latest, err := s.GetLatestChecks(ctx)
	if err != nil {
		t.Fatalf("get latest checks: %v", err)
	}
	if len(latest) != 1 {
		t.Fatalf("expected 1 latest check, got %d", len(latest))
	}
}
