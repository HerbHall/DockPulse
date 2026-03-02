package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	socketPath := flag.String("socket", "/run/guest-services/dockpulse.sock", "Unix socket path")
	dbPath := flag.String("db", "/data/dockpulse.db", "SQLite database path")
	flag.Parse()

	log.Printf("DockPulse starting: socket=%s db=%s", *socketPath, *dbPath)

	// TODO: Initialize store, checker, registry client, HTTP handlers
	// TODO: Listen on Unix socket

	log.Println("DockPulse placeholder -- no server yet")
	os.Exit(0)
}
