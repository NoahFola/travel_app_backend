package handlers

import (
	"net/http"

	"github.com/NoahFola/travel_app_backend/internal/repository"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	DeviceRepo *repository.DeviceRepository
}

type registerDeviceRequest struct {
	Token string `json:"token" binding:"required"`
}

func (h *UserHandler) RegisterDevice(c *gin.Context) {
	var req registerDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.DeviceRepo.RegisterToken(c.Request.Context(), userID, req.Token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "device registered"})
}
