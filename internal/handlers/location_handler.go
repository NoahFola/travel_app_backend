package handlers

import (
	"net/http"

	"github.com/NoahFola/travel_app_backend/internal/service"
	"github.com/gin-gonic/gin"
)

type LocationHandler struct {
	Service *service.LocationService
}

func (h *LocationHandler) Search(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
		return
	}

	results, err := h.Service.SearchPlaces(query)
	if err != nil {
		// Log error internally
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
