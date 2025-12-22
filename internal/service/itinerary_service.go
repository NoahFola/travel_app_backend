package service

import (
	"context"
	"errors"

	"github.com/NoahFola/travel_app_backend/internal/domain"
	"github.com/NoahFola/travel_app_backend/internal/repository"
)

type ItineraryService struct {
	Repo     *repository.ItineraryRepository
	TripRepo *repository.TripRepository
}

func NewItineraryService(repo *repository.ItineraryRepository, tripRepo *repository.TripRepository) *ItineraryService {
	return &ItineraryService{
		Repo:     repo,
		TripRepo: tripRepo,
	}
}

func (s *ItineraryService) CreateItinerary(ctx context.Context, itinerary *domain.Itinerary) error {
	// Verify trip exists
	_, err := s.TripRepo.GetByID(ctx, itinerary.TripID)
	if err != nil {
		return errors.New("trip not found")
	}
	return s.Repo.Create(ctx, itinerary)
}

func (s *ItineraryService) GetItinerary(ctx context.Context, id string) (*domain.Itinerary, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *ItineraryService) ListItineraries(ctx context.Context, tripID string) ([]domain.Itinerary, error) {
	// Verify trip exists
	_, err := s.TripRepo.GetByID(ctx, tripID)
	if err != nil {
		return nil, errors.New("trip not found")
	}
	return s.Repo.GetByTripID(ctx, tripID)
}

func (s *ItineraryService) UpdateItinerary(ctx context.Context, itinerary *domain.Itinerary) error {
	return s.Repo.Update(ctx, itinerary)
}

func (s *ItineraryService) DeleteItinerary(ctx context.Context, id string) error {
	return s.Repo.Delete(ctx, id)
}
