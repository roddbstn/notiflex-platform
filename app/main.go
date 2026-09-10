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

const version = "v0.7.0"

const valkeyAddrFile = "/mnt/secrets/valkey-addr"

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
	data, err := os.ReadFile(valkeyAddrFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read valkey addr: %v", err), http.StatusInternalServerError)
		return
	}
	addr := strings.TrimSpace(string(data))

	password := valkeyPassword()

	id, err := valkeyIncr(addr, password, "notiflex:counter")
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

// valkeyPassword returns the Valkey password.
// If VALKEY_PASSWORD_FILE is set, reads from that file.
// Falls back to VALKEY_PASSWORD env var.
func valkeyPassword() string {
	if path := os.Getenv("VALKEY_PASSWORD_FILE"); path != "" {
		data, err := os.ReadFile(path)
		if err == nil {
			return strings.TrimSpace(string(data))
		}
	}
	return os.Getenv("VALKEY_PASSWORD")
}

// valkeyIncr sends AUTH (if password set) then INCR via RESP protocol.
func valkeyIncr(addr, password, key string) (int64, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	if password != "" {
		auth := fmt.Sprintf("*2\r\n$4\r\nAUTH\r\n$%d\r\n%s\r\n", len(password), password)
		if _, err := fmt.Fprint(conn, auth); err != nil {
			return 0, err
		}
		line, err := reader.ReadString('\n')
		if err != nil {
			return 0, err
		}
		if !strings.HasPrefix(strings.TrimSpace(line), "+OK") {
			return 0, fmt.Errorf("AUTH failed: %s", strings.TrimSpace(line))
		}
	}

	cmd := fmt.Sprintf("*2\r\n$4\r\nINCR\r\n$%d\r\n%s\r\n", len(key), key)
	if _, err := fmt.Fprint(conn, cmd); err != nil {
		return 0, err
	}

	line, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	line = strings.TrimSpace(line)
	if line[0] != ':' {
		return 0, fmt.Errorf("unexpected RESP response: %s", line)
	}

	return strconv.ParseInt(line[1:], 10, 64)
}
