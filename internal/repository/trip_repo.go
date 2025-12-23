package repository

import (
	"context"
	"errors"

	"github.com/NoahFola/travel_app_backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TripRepository struct {
	DB *pgxpool.Pool
}

func NewTripRepository(db *pgxpool.Pool) *TripRepository {
	return &TripRepository{DB: db}
}

func (r *TripRepository) Create(ctx context.Context, trip *domain.Trip) error {
	query := `
		INSERT INTO trips (user_id, location, start_date, end_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	err := r.DB.QueryRow(ctx, query, trip.UserID, trip.Location, trip.StartDate, trip.EndDate).Scan(&trip.ID, &trip.CreatedAt, &trip.UpdatedAt)
	print("Trip created: ", trip.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *TripRepository) GetByID(ctx context.Context, id string) (*domain.Trip, error) {
	query := `
		SELECT id, user_id, location, start_date, end_date, created_at, updated_at
		FROM trips
		WHERE id = $1`

	var trip domain.Trip
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&trip.ID,
		&trip.UserID,
		&trip.Location,
		&trip.StartDate,
		&trip.EndDate,
		&trip.CreatedAt,
		&trip.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("trip not found")
		}
		return nil, err
	}
	return &trip, nil
}

func (r *TripRepository) GetByUserID(ctx context.Context, userID string) ([]domain.Trip, error) {
	query := `
		SELECT id, user_id, location, start_date, end_date, created_at, updated_at
		FROM trips
		WHERE user_id = $1
		ORDER BY start_date DESC`

	rows, err := r.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trips []domain.Trip
	for rows.Next() {
		var trip domain.Trip
		err := rows.Scan(
			&trip.ID,
			&trip.UserID,
			&trip.Location,
			&trip.StartDate,
			&trip.EndDate,
			&trip.CreatedAt,
			&trip.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		trips = append(trips, trip)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return trips, nil
}

func (r *TripRepository) Update(ctx context.Context, trip *domain.Trip) error {
	query := `
		UPDATE trips
		SET location = $1, start_date = $2, end_date = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING updated_at`

	err := r.DB.QueryRow(ctx, query, trip.Location, trip.StartDate, trip.EndDate, trip.ID).Scan(&trip.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("trip not found")
		}
		return err
	}
	return nil
}

func (r *TripRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM trips WHERE id = $1`

	commandTag, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return errors.New("trip not found")
	}
	return nil
}
