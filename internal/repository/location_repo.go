package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Location struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Address       string  `json:"address"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	GooglePlaceID *string `json:"google_place_id"`
}

type LocationRepository struct {
	DB *pgxpool.Pool
}

func NewLocationRepository(db *pgxpool.Pool) *LocationRepository {
	return &LocationRepository{DB: db}
}

func (r *LocationRepository) Create(ctx context.Context, loc *Location) error {
	query := `
		INSERT INTO locations (name, address, latitude, longitude, google_place_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err := r.DB.QueryRow(ctx, query, loc.Name, loc.Address, loc.Latitude, loc.Longitude, loc.GooglePlaceID).Scan(&loc.ID)
	return err
}

func (r *LocationRepository) GetByPlaceID(ctx context.Context, placeID string) (*Location, error) {
	query := `
		SELECT id, name, address, latitude, longitude, google_place_id
		FROM locations
		WHERE google_place_id = $1
	`
	var loc Location
	err := r.DB.QueryRow(ctx, query, placeID).Scan(
		&loc.ID, &loc.Name, &loc.Address, &loc.Latitude, &loc.Longitude, &loc.GooglePlaceID,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &loc, nil
}

func (r *LocationRepository) GetByID(ctx context.Context, id string) (*Location, error) {
	query := `
		SELECT id, name, address, latitude, longitude, google_place_id
		FROM locations
		WHERE id = $1
	`
	var loc Location
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&loc.ID, &loc.Name, &loc.Address, &loc.Latitude, &loc.Longitude, &loc.GooglePlaceID,
	)
	if err == pgx.ErrNoRows {
		return nil, nil // Return nil if not found
	}
	if err != nil {
		return nil, err
	}
	return &loc, nil
}
