package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
)

const version = "v0.1.3"

var counter atomic.Int64

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/id", idHandler)
	http.HandleFunc("/version", versionHandler)
	fmt.Println("Notiflex API server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"version": version})
}

func idHandler(w http.ResponseWriter, r *http.Request) {
	id := counter.Add(1)
	pod := os.Getenv("POD_NAME")
	if pod == "" {
		pod = "unknown"
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"id":  id,
		"pod": pod,
	})
}
