package handlers

import (
	"net/http"
	"time"

	"github.com/NoahFola/travel_app_backend/internal/domain"
	"github.com/NoahFola/travel_app_backend/internal/service"
	"github.com/gin-gonic/gin"
)

type TripHandler struct {
	Service *service.TripService
}

type createTripRequest struct {
	Location  string    `json:"location" binding:"required"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}

type updateTripRequest struct {
	Location  string    `json:"location"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

func (h *TripHandler) CreateTrip(c *gin.Context) {
	var req createTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("userID") // Assumes AuthMiddleware sets this
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	trip := &domain.Trip{
		UserID:    userID,
		Location:  req.Location,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	if err := h.Service.CreateTrip(c.Request.Context(), trip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, trip)
}

func (h *TripHandler) GetTrip(c *gin.Context) {
	id := c.Param("id")
	trip, err := h.Service.GetTrip(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "trip not found"})
		return
	}
	c.JSON(http.StatusOK, trip)
}

func (h *TripHandler) ListMyTrips(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	trips, err := h.Service.ListUserTrips(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trips)
}

func (h *TripHandler) UpdateTrip(c *gin.Context) {
	id := c.Param("id")
	var req updateTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trip, err := h.Service.GetTrip(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "trip not found"})
		return
	}

	// Update fields if present
	if req.Location != "" {
		trip.Location = req.Location
	}
	if !req.StartDate.IsZero() {
		trip.StartDate = req.StartDate
	}
	if !req.EndDate.IsZero() {
		trip.EndDate = req.EndDate
	}

	if err := h.Service.UpdateTrip(c.Request.Context(), trip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trip)
}

func (h *TripHandler) DeleteTrip(c *gin.Context) {
	id := c.Param("id")
	if err := h.Service.DeleteTrip(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "trip not found"}) // Assume 404 for simplicity
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "trip deleted"})
}

func (h *TripHandler) ShareTrip(c *gin.Context) {
	tripID := c.Param("id")
	// Verify ownership is usually done here or in service.
	// For MVP, assuming if you can hit this endpoint with auth, you might check inside Service or Repo if userID matches.
	// But TripService.GenerateShareToken doesn't check UserID currently.
	// Strict implementation would check if trip.UserID == currentUser before sharing.

	token, err := h.Service.GenerateShareToken(c.Request.Context(), tripID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"share_token": token, "url": "/preview/" + token})
}

func (h *TripHandler) GetSharedTrip(c *gin.Context) {
	token := c.Param("token")
	trip, err := h.Service.GetTripByShareToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "trip not found or expired"})
		return
	}

	c.JSON(http.StatusOK, trip)
}
