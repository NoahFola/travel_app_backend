package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Media struct {
	ID         string `json:"id"`
	URL        string `json:"url"`
	Type       string `json:"type"`
	ActivityID string `json:"activity_id"`
}

type MediaRepository struct {
	DB *pgxpool.Pool
}

func NewMediaRepository(db *pgxpool.Pool) *MediaRepository {
	return &MediaRepository{DB: db}
}

func (r *MediaRepository) Create(ctx context.Context, media *Media) error {
	query := `
		INSERT INTO media (url, type, activity_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	err := r.DB.QueryRow(ctx, query, media.URL, media.Type, media.ActivityID).Scan(&media.ID)
	return err
}

func (r *MediaRepository) ListByActivityID(ctx context.Context, activityID string) ([]Media, error) {
	query := `
		SELECT id, url, type, activity_id
		FROM media
		WHERE activity_id = $1
	`
	rows, err := r.DB.Query(ctx, query, activityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var medias []Media
	for rows.Next() {
		var m Media
		if err := rows.Scan(&m.ID, &m.URL, &m.Type, &m.ActivityID); err != nil {
			return nil, err
		}
		medias = append(medias, m)
	}
	return medias, nil
}
