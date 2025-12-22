package handlers

import (
	"net/http"

	"github.com/NoahFola/travel_app_backend/internal/service"
	"github.com/gin-gonic/gin"
)

type MediaHandler struct {
	Service *service.MediaService
}

func (h *MediaHandler) Upload(c *gin.Context) {
	// 1. Get Activity ID
	activityID := c.PostForm("activity_id")
	if activityID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "activity_id is required"})
		return
	}

	// 2. Get File
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	// 3. Upload
	media, err := h.Service.UploadMedia(c.Request.Context(), file, activityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, media)
}
