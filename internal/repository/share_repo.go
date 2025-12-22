package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ShareToken struct {
	ID        string    `json:"id"`
	TripID    string    `json:"trip_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type ShareRepository struct {
	DB *pgxpool.Pool
}

func NewShareRepository(db *pgxpool.Pool) *ShareRepository {
	return &ShareRepository{DB: db}
}

func (r *ShareRepository) CreateToken(ctx context.Context, tripID, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO share_tokens (trip_id, token, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var id string
	err := r.DB.QueryRow(ctx, query, tripID, token, expiresAt).Scan(&id)
	return err
}

func (r *ShareRepository) GetTripIDByToken(ctx context.Context, token string) (string, error) {
	query := `
		SELECT trip_id
		FROM share_tokens
		WHERE token = $1 AND expires_at > NOW()
	`
	var tripID string
	err := r.DB.QueryRow(ctx, query, token).Scan(&tripID)
	if err == pgx.ErrNoRows {
		return "", nil // Token not found or expired
	}
	if err != nil {
		return "", err
	}
	return tripID, nil
}
