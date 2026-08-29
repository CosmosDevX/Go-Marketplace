// Package handler
package handler

import (
	"myapp/internal/logger"
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

// CheckHealth godoc
//
//	@Summary		Health check
//	@Description	Check availability of database and Redis.
//	@Tags			Health
//	@Produce		json
//	@Success		200	{object}	StatusResponse	"Service is healthy"
//	@Failure		503	{object}	StatusResponse	"Service is unhealthy"
//	@Router			/health [get]
func (h HealthHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	if err := h.db.PingContext(ctx); err != nil {
		log.Error("sql db is unavailable", "error", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		utils.WriteJSON(w, map[string]string{"status": "unhealthy"})
		return
	}

	if err := h.redisClient.Ping(ctx).Err(); err != nil {
		log.Error("redis is unavailable", "error", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		utils.WriteJSON(w, map[string]string{"status": "unhealthy"})
		return
	}

	utils.WriteJSON(w, map[string]string{"status": "ok"})
}
