package service

import (
	"context"
	"errors"

	"github.com/NoahFola/travel_app_backend/internal/domain"
	"github.com/NoahFola/travel_app_backend/internal/repository"
)

type ActivityService struct {
	Repo          *repository.ActivityRepository
	ItineraryRepo *repository.ItineraryRepository
}

func NewActivityService(repo *repository.ActivityRepository, itineraryRepo *repository.ItineraryRepository) *ActivityService {
	return &ActivityService{
		Repo:          repo,
		ItineraryRepo: itineraryRepo,
	}
}

func (s *ActivityService) CreateActivity(ctx context.Context, activity *domain.Activity) error {
	// Verify itinerary exists if provided (it MUST be provided for now based on handler)
	if activity.ItineraryID == nil {
		return errors.New("itinerary_id is required")
	}

	itinerary, err := s.ItineraryRepo.GetByID(ctx, *activity.ItineraryID)
	if err != nil {
		return errors.New("itinerary not found")
	}

	// Set the TripID from the Itinerary
	activity.TripID = itinerary.TripID

	return s.Repo.Create(ctx, activity)
}

func (s *ActivityService) GetActivity(ctx context.Context, id string) (*domain.Activity, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *ActivityService) ListActivities(ctx context.Context, itineraryID string) ([]domain.Activity, error) {
	// Verify itinerary exists
	_, err := s.ItineraryRepo.GetByID(ctx, itineraryID)
	if err != nil {
		return nil, errors.New("itinerary not found")
	}
	return s.Repo.GetByItineraryID(ctx, itineraryID)
}

func (s *ActivityService) UpdateActivity(ctx context.Context, activity *domain.Activity) error {
	// If itinerary ID changed, we might want to verify it exists and belongs to same trip?
	// For now, simple update.
	return s.Repo.Update(ctx, activity)
}

func (s *ActivityService) DeleteActivity(ctx context.Context, id string) error {
	return s.Repo.Delete(ctx, id)
}
