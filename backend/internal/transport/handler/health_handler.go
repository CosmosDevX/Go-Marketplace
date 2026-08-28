// Package handler
package handler

import (
	"log"
	"myapp/internal/utils"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type HealthHandler struct {
	db          *sqlx.DB
	redisClient *redis.Client
}

func NewHealthHandler(db *sqlx.DB, redisClient *redis.Client) HealthHandler {
	return HealthHandler{
		db:          db,
		redisClient: redisClient,
	}
}

func (h HealthHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := h.db.PingContext(ctx); err != nil {
		log.Println("sql db is unvailable")
		w.WriteHeader(http.StatusServiceUnavailable)
		utils.WriteJSON(w, map[string]string{"status": "unhealthy"})
		return
	}

	if err := h.redisClient.Ping(ctx).Err(); err != nil {
		log.Println("redis db is unvailable")
		w.WriteHeader(http.StatusServiceUnavailable)
		utils.WriteJSON(w, map[string]string{"status": "unhealthy"})
		return
	}

	utils.WriteJSON(w, map[string]string{"status": "ok"})
}
