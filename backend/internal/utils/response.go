package utils

import (
	"context"
	"encoding/json"
	"net"
	"net/http"

	"myapp/internal/logger"
)

func WriteJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "error during writing response", http.StatusInternalServerError)
		return
	}
}

func WriteError(ctx context.Context, w http.ResponseWriter, err error) {
	status := MapError(err)
	msg := PublicMessage(err)

	logger.FromContext(ctx).Error("request failed", "status", status, "error", err)

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
