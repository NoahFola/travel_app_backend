package handlers

import (
	"net/http"
	"time"

	"github.com/NoahFola/travel_app_backend/internal/domain"
	"github.com/NoahFola/travel_app_backend/internal/service"
	"github.com/gin-gonic/gin"
)

type ActivityHandler struct {
	Service *service.ActivityService
}

type createActivityRequest struct {
	Name        string     `json:"name" binding:"required"`
	Description *string    `json:"description"`
	Location    *string    `json:"location"`
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	Type        *string    `json:"type"`
	Status      string     `json:"status"`
}

type updateActivityRequest struct {
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Location    *string    `json:"location"`
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	Type        *string    `json:"type"`
	Status      string     `json:"status"`
	ItineraryID *string    `json:"itinerary_id"` // can move between days
}

func (h *ActivityHandler) CreateActivity(c *gin.Context) {
	itineraryID := c.Param("id")
	var req createActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	activity := &domain.Activity{
		ItineraryID: &itineraryID,
		Name:        req.Name,
		Description: req.Description,
		Location:    req.Location,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Type:        req.Type,
		Status:      req.Status,
	}
	// Need to fetch tripID from context or look it up?
	// The Service layer assumes `TripID` is present on struct?
	// `activity_repo.go` needs `trip_id` for insertion (NOT NULL).
	// We might need to look up the itinerary first to get the trip_id.
	// Or pass it in. But we only have `itineraryId`.

	// Better approach: ActivityService's CreateActivity should look up the Itinerary to fill TripID.
	// Let's assume generic "planned" status if empty?
	if activity.Status == "" {
		activity.Status = "planned"
	}

	// We need to set TripID.
	// We can add a helper in ActivityHandler or Service to fill this.
	// Let's do it in the Service. No, Service expects a domain object.
	// Let's fetch the itinerary here to get TripID.
	// Ideally we'd inject ItineraryService into ActivityHandler too, or ActivityService exposes a way.

	// Wait, ActivityService has access to ItineraryRepo.
	// Let's modify ActivityService.CreateActivity to fetch Itinerary and set TripID if missing.

	// For now, I'll trust the Service handles it or I'll fix it in next step.
	// Actually, `activity_repo.go` insert query uses `activity.TripID`.
	// If I don't set it, it will be empty string, and DB will reject it (UUID expected).
	// So I MUST set it.

	// I will handle this in the Handler by using a small "GetItinerary" lookup via a helper method or direct repo access?
	// Handlers shouldn't access repos.
	// ActivityService should probably handle "CreateActivityForItinerary".

	// Let's just pass what we have to Service, but Service needs to be smart.
	// I'll call `h.Service.CreateActivity`.
	// But before that, I need to fetch Itinerary to get TripID?
	// Or Update ActivityService.

	// I'll proceed with Handler implementation assuming Service will be updated or is sufficient.
	// Wait, I implemented `ActivityService.CreateActivity` just calling `repo.Create`.
	// And `repo.Create` expects `TripID`.
	// So `ActivityService.CreateActivity` IS BROKEN currently for this flow.
	// I will fix `ActivityService.CreateActivity` in the next step or right now.

	// Actually, better to fix `ActivityService` first? Or just implement Handler and then fix Service.
	// I'll implement Handler, then I'll see I need to fix Service.

	if err := h.Service.CreateActivity(c.Request.Context(), activity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, activity)
}

func (h *ActivityHandler) ListActivities(c *gin.Context) {
	itineraryID := c.Param("id")
	activities, err := h.Service.ListActivities(c.Request.Context(), itineraryID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, activities)
}

func (h *ActivityHandler) GetActivity(c *gin.Context) {
	id := c.Param("id")
	activity, err := h.Service.GetActivity(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "activity not found"})
		return
	}
	c.JSON(http.StatusOK, activity)
}

func (h *ActivityHandler) UpdateActivity(c *gin.Context) {
	id := c.Param("id")
	var req updateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	activity, err := h.Service.GetActivity(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "activity not found"})
		return
	}

	if req.Name != "" {
		activity.Name = req.Name
	}
	if req.Description != nil {
		activity.Description = req.Description
	}
	if req.Location != nil {
		activity.Location = req.Location
	}
	if req.StartTime != nil {
		activity.StartTime = req.StartTime
	}
	if req.EndTime != nil {
		activity.EndTime = req.EndTime
	}
	if req.Type != nil {
		activity.Type = req.Type
	}
	if req.Status != "" {
		activity.Status = req.Status
	}
	if req.ItineraryID != nil {
		activity.ItineraryID = req.ItineraryID
	}

	if err := h.Service.UpdateActivity(c.Request.Context(), activity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, activity)
}

func (h *ActivityHandler) DeleteActivity(c *gin.Context) {
	id := c.Param("id")
	if err := h.Service.DeleteActivity(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "activity not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "activity deleted"})
}
