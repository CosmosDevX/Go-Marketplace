package utils

import (
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "error during writing response", http.StatusInternalServerError)
		return
	}
}

func WriteError(w http.ResponseWriter, err error) {
	status := MapError(err)
	msg := PublicMessage(err)

	slog.Error("request failed", "status", status, "error", err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": msg})
}

func GetIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return ip
}
