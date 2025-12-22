package repository

import (
	"context"
	"errors"

	"github.com/NoahFola/travel_app_backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ActivityRepository struct {
	DB *pgxpool.Pool
}

func NewActivityRepository(db *pgxpool.Pool) *ActivityRepository {
	return &ActivityRepository{DB: db}
}

func (r *ActivityRepository) Create(ctx context.Context, activity *domain.Activity) error {
	query := `
		INSERT INTO activities (trip_id, itinerary_id, name, description, location, start_time, end_time, type, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	err := r.DB.QueryRow(ctx, query,
		activity.TripID,
		activity.ItineraryID,
		activity.Name,
		activity.Description,
		activity.Location,
		activity.StartTime,
		activity.EndTime,
		activity.Type,
		activity.Status,
	).Scan(&activity.ID, &activity.CreatedAt, &activity.UpdatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (r *ActivityRepository) GetByID(ctx context.Context, id string) (*domain.Activity, error) {
	query := `
		SELECT id, trip_id, itinerary_id, name, description, location, start_time, end_time, type, status, created_at, updated_at
		FROM activities
		WHERE id = $1`

	var a domain.Activity
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&a.ID,
		&a.TripID,
		&a.ItineraryID,
		&a.Name,
		&a.Description,
		&a.Location,
		&a.StartTime,
		&a.EndTime,
		&a.Type,
		&a.Status,
		&a.CreatedAt,
		&a.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("activity not found")
		}
		return nil, err
	}
	return &a, nil
}

func (r *ActivityRepository) GetByItineraryID(ctx context.Context, itineraryID string) ([]domain.Activity, error) {
	query := `
		SELECT id, trip_id, itinerary_id, name, description, location, start_time, end_time, type, status, created_at, updated_at
		FROM activities
		WHERE itinerary_id = $1
		ORDER BY start_time ASC`

	rows, err := r.DB.Query(ctx, query, itineraryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []domain.Activity
	for rows.Next() {
		var a domain.Activity
		err := rows.Scan(
			&a.ID,
			&a.TripID,
			&a.ItineraryID,
			&a.Name,
			&a.Description,
			&a.Location,
			&a.StartTime,
			&a.EndTime,
			&a.Type,
			&a.Status,
			&a.CreatedAt,
			&a.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		activities = append(activities, a)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return activities, nil
}

func (r *ActivityRepository) Update(ctx context.Context, activity *domain.Activity) error {
	query := `
		UPDATE activities
		SET itinerary_id = $1, name = $2, description = $3, location = $4, start_time = $5, end_time = $6, type = $7, status = $8, updated_at = NOW()
		WHERE id = $9
		RETURNING updated_at`

	err := r.DB.QueryRow(ctx, query,
		activity.ItineraryID,
		activity.Name,
		activity.Description,
		activity.Location,
		activity.StartTime,
		activity.EndTime,
		activity.Type,
		activity.Status,
		activity.ID,
	).Scan(&activity.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("activity not found")
		}
		return err
	}
	return nil
}

func (r *ActivityRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM activities WHERE id = $1`
	ct, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return errors.New("activity not found")
	}
	return nil
}
