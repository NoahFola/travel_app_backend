package service

import (
	"context"

	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/NoahFola/travel_app_backend/internal/domain"
	"github.com/NoahFola/travel_app_backend/internal/repository"
)

type TripService struct {
	Repo      *repository.TripRepository
	ShareRepo *repository.ShareRepository
}

func NewTripService(repo *repository.TripRepository, shareRepo *repository.ShareRepository) *TripService {
	return &TripService{Repo: repo, ShareRepo: shareRepo}
}

func (s *TripService) CreateTrip(ctx context.Context, trip *domain.Trip) error {
	return s.Repo.Create(ctx, trip)
}

func (s *TripService) GetTrip(ctx context.Context, id string) (*domain.Trip, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *TripService) ListUserTrips(ctx context.Context, userID string) ([]domain.Trip, error) {
	return s.Repo.GetByUserID(ctx, userID)
}

func (s *TripService) UpdateTrip(ctx context.Context, trip *domain.Trip) error {
	return s.Repo.Update(ctx, trip)
}

func (s *TripService) DeleteTrip(ctx context.Context, id string) error {
	return s.Repo.Delete(ctx, id)
}

// Share Logic

func (s *TripService) GenerateShareToken(ctx context.Context, tripID string) (string, error) {
	// Generate random token
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	token := hex.EncodeToString(bytes)

	// Expires in 30 days
	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	if err := s.ShareRepo.CreateToken(ctx, tripID, token, expiresAt); err != nil {
		return "", err
	}

	return token, nil
}

func (s *TripService) GetTripByShareToken(ctx context.Context, token string) (*domain.Trip, error) {
	tripID, err := s.ShareRepo.GetTripIDByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if tripID == "" {
		return nil, errors.New("invalid or crossed share token")
	}

	return s.Repo.GetByID(ctx, tripID)
}
