package handlers

import (
	"net/http"
	"time"

	"github.com/NoahFola/travel_app_backend/internal/domain"
	"github.com/NoahFola/travel_app_backend/internal/service"
	"github.com/gin-gonic/gin"
)

type ItineraryHandler struct {
	Service *service.ItineraryService
}

type createItineraryRequest struct {
	Slug  string    `json:"slug" binding:"required"`
	Title *string   `json:"title"`
	Date  time.Time `json:"date" binding:"required"`
}

type updateItineraryRequest struct {
	Slug  string    `json:"slug"`
	Title *string   `json:"title"`
	Date  time.Time `json:"date"`
}

func (h *ItineraryHandler) CreateItinerary(c *gin.Context) {
	tripID := c.Param("tripId")
	var req createItineraryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	itinerary := &domain.Itinerary{
		TripID: tripID,
		Slug:   req.Slug,
		Title:  req.Title,
		Date:   req.Date,
	}

	if err := h.Service.CreateItinerary(c.Request.Context(), itinerary); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, itinerary)
}

func (h *ItineraryHandler) ListItineraries(c *gin.Context) {
	tripID := c.Param("tripId")
	itineraries, err := h.Service.ListItineraries(c.Request.Context(), tripID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, itineraries)
}

func (h *ItineraryHandler) GetItinerary(c *gin.Context) {
	id := c.Param("id")
	itinerary, err := h.Service.GetItinerary(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "itinerary not found"})
		return
	}
	c.JSON(http.StatusOK, itinerary)
}

func (h *ItineraryHandler) UpdateItinerary(c *gin.Context) {
	id := c.Param("id")
	var req updateItineraryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	itinerary, err := h.Service.GetItinerary(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "itinerary not found"})
		return
	}

	if req.Slug != "" {
		itinerary.Slug = req.Slug
	}
	if req.Title != nil {
		itinerary.Title = req.Title
	}
	if !req.Date.IsZero() {
		itinerary.Date = req.Date
	}

	if err := h.Service.UpdateItinerary(c.Request.Context(), itinerary); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, itinerary)
}

func (h *ItineraryHandler) DeleteItinerary(c *gin.Context) {
	id := c.Param("id")
	if err := h.Service.DeleteItinerary(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "itinerary not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "itinerary deleted"})
}
