package repository

import (
	"context"
	"errors"

	"github.com/NoahFola/travel_app_backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ItineraryRepository struct {
	DB *pgxpool.Pool
}

func NewItineraryRepository(db *pgxpool.Pool) *ItineraryRepository {
	return &ItineraryRepository{DB: db}
}

func (r *ItineraryRepository) Create(ctx context.Context, itinerary *domain.Itinerary) error {
	query := `
		INSERT INTO itineraries (trip_id, slug, title, date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	err := r.DB.QueryRow(ctx, query, itinerary.TripID, itinerary.Slug, itinerary.Title, itinerary.Date).
		Scan(&itinerary.ID, &itinerary.CreatedAt, &itinerary.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *ItineraryRepository) GetByID(ctx context.Context, id string) (*domain.Itinerary, error) {
	query := `
		SELECT id, trip_id, slug, title, date, created_at, updated_at
		FROM itineraries
		WHERE id = $1`

	var i domain.Itinerary
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&i.ID,
		&i.TripID,
		&i.Slug,
		&i.Title,
		&i.Date,
		&i.CreatedAt,
		&i.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("itinerary not found")
		}
		return nil, err
	}
	return &i, nil
}

func (r *ItineraryRepository) GetByTripID(ctx context.Context, tripID string) ([]domain.Itinerary, error) {
	query := `
		SELECT id, trip_id, slug, title, date, created_at, updated_at
		FROM itineraries
		WHERE trip_id = $1
		ORDER BY date ASC`

	rows, err := r.DB.Query(ctx, query, tripID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var itineraries []domain.Itinerary
	for rows.Next() {
		var i domain.Itinerary
		err := rows.Scan(
			&i.ID,
			&i.TripID,
			&i.Slug,
			&i.Title,
			&i.Date,
			&i.CreatedAt,
			&i.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		itineraries = append(itineraries, i)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return itineraries, nil
}

func (r *ItineraryRepository) Update(ctx context.Context, itinerary *domain.Itinerary) error {
	query := `
		UPDATE itineraries
		SET slug = $1, title = $2, date = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING updated_at`

	err := r.DB.QueryRow(ctx, query, itinerary.Slug, itinerary.Title, itinerary.Date, itinerary.ID).
		Scan(&itinerary.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("itinerary not found")
		}
		return err
	}
	return nil
}

func (r *ItineraryRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM itineraries WHERE id = $1`
	ct, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return errors.New("itinerary not found")
	}
	return nil
}
