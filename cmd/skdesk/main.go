package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/domehahn/skdesk/pkg/server"
	"github.com/domehahn/skdesk/pkg/signer"
	"github.com/domehahn/skdesk/pkg/ticket"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "version":
		fmt.Printf("skdesk version %s\n", version)
	case "keygen":
		handleKeygen()
	case "serve":
		handleServe(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`skdesk — Human-in-the-Loop Review & Attestation Sign-off Portal

Usage:
  skdesk keygen
  skdesk serve [--port 8080]
  skdesk version`)
}

func handleKeygen() {
	kp, err := signer.GenerateKeyPair()
	if err != nil {
		fmt.Fprintf(os.Stderr, "keygen error: %v\n", err)
		os.Exit(1)
	}
	out, _ := json.MarshalIndent(kp, "", "  ")
	fmt.Println(string(out))
}

func handleServe(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	port := fs.Int("port", 8080, "HTTP server port")
	_ = fs.Parse(args)

	mgr := ticket.NewManager()
	// Seed demo ticket
	mgr.Create("database-migrator", "sha256:7a8b9c...", "production", "Requires external database write access")

	srv := server.NewServer(mgr)
	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("skdesk server listening on http://localhost%s\n", addr)
	_ = http.ListenAndServe(addr, srv.Routes())
}
