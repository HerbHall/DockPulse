package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/HerbHall/DockPulse/backend/internal/checker"
	"github.com/HerbHall/DockPulse/backend/internal/store"
)

// Handler provides HTTP endpoints for the DockPulse backend API.
type Handler struct {
	checker *checker.Checker
	store   *store.Store
}

// NewHandler creates a Handler wired to the given checker and store.
func NewHandler(c *checker.Checker, s *store.Store) *Handler {
	return &Handler{
		checker: c,
		store:   s,
	}
}

// RegisterRoutes mounts the API endpoints on the given ServeMux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/checks", h.getChecks)
	mux.HandleFunc("POST /api/check-all", h.checkAll)
	mux.HandleFunc("GET /api/status", h.getStatus)
}

// checksResponse wraps a list of image checks for JSON serialization.
type checksResponse struct {
	Checks []store.ImageCheck `json:"checks"`
}

// checkAllResponse includes both the check results and when the check started.
type checkAllResponse struct {
	Checks    []store.ImageCheck `json:"checks"`
	StartedAt string             `json:"startedAt"`
}

// statusResponse represents the health check response.
type statusResponse struct {
	Healthy bool   `json:"healthy"`
	Version string `json:"version"`
}

func (h *Handler) getChecks(w http.ResponseWriter, r *http.Request) {
	checks, err := h.store.GetLatestChecks(r.Context())
	if err != nil {
		log.Printf("api: get checks: %v", err)
		writeError(w, http.StatusInternalServerError)
		return
	}

	if checks == nil {
		checks = []store.ImageCheck{}
	}

	writeJSON(w, http.StatusOK, checksResponse{Checks: checks})
}

func (h *Handler) checkAll(w http.ResponseWriter, r *http.Request) {
	startedAt := time.Now().UTC().Format(time.RFC3339)

	results, err := h.checker.CheckAll(r.Context())
	if err != nil {
		log.Printf("api: check all: %v", err)
		writeError(w, http.StatusInternalServerError)
		return
	}

	if results == nil {
		results = []store.ImageCheck{}
	}

	writeJSON(w, http.StatusOK, checkAllResponse{
		Checks:    results,
		StartedAt: startedAt,
	})
}

func (h *Handler) getStatus(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, statusResponse{
		Healthy: true,
		Version: "0.1.0",
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("api: encode json: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": http.StatusText(status)}); err != nil {
		log.Printf("api: encode error: %v", err)
	}
}
