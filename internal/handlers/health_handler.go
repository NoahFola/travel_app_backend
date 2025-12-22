package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
	DB *pgxpool.Pool
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	if err := h.DB.Ping(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "db": "disconnected"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "db": "connected"})
}
