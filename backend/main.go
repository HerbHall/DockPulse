package main

import (
	"errors"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HerbHall/DockPulse/backend/internal/api"
	"github.com/HerbHall/DockPulse/backend/internal/checker"
	"github.com/HerbHall/DockPulse/backend/internal/docker"
	"github.com/HerbHall/DockPulse/backend/internal/registry"
	"github.com/HerbHall/DockPulse/backend/internal/store"
)

func main() {
	socketPath := flag.String("socket", "/run/guest-services/dockpulse.sock", "Unix socket path")
	dbPath := flag.String("db", "/data/dockpulse.db", "SQLite database path")
	flag.Parse()

	log.Printf("DockPulse starting: socket=%s db=%s", *socketPath, *dbPath)

	if err := run(*socketPath, *dbPath); err != nil {
		log.Fatalf("Fatal: %v", err)
	}
}

func run(socketPath, dbPath string) error {
	// Initialize store.
	s, err := store.New(dbPath)
	if err != nil {
		return err
	}
	defer func() { _ = s.Close() }()

	// Initialize Docker client.
	dc, err := docker.NewClient()
	if err != nil {
		return err
	}
	defer func() { _ = dc.Close() }()

	// Initialize registry client.
	reg := registry.NewDockerHubClient()

	// Initialize checker.
	chk := checker.New(s, reg, dc)

	// Set up HTTP handler.
	handler := api.NewHandler(chk, s)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	// Remove stale socket file if it exists.
	_ = os.RemoveAll(socketPath)

	ln, err := net.Listen("unix", socketPath)
	if err != nil {
		return err
	}

	log.Printf("DockPulse listening on %s", socketPath)

	srv := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Graceful shutdown on SIGTERM/SIGINT.
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
		<-sigCh
		log.Println("Shutting down...")
		_ = ln.Close()
	}()

	if err = srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) && !errors.Is(err, net.ErrClosed) {
		return err
	}

	return nil
}
