package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const version = "v0.4.0"

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
	addr := os.Getenv("VALKEY_ADDR")
	if addr == "" {
		addr = "valkey-primary:6379"
	}

	id, err := valkeyIncr(addr, "notiflex:counter")
	if err != nil {
		http.Error(w, fmt.Sprintf("valkey error: %v", err), http.StatusInternalServerError)
		return
	}

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

// valkeyIncr sends INCR key via RESP protocol and returns the new value.
func valkeyIncr(addr, key string) (int64, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	cmd := fmt.Sprintf("*2\r\n$4\r\nINCR\r\n$%d\r\n%s\r\n", len(key), key)
	if _, err := fmt.Fprint(conn, cmd); err != nil {
		return 0, err
	}

	line, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return 0, err
	}

	line = strings.TrimSpace(line)
	if line[0] != ':' {
		return 0, fmt.Errorf("unexpected RESP response: %s", line)
	}

	return strconv.ParseInt(line[1:], 10, 64)
}
