package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DeviceRepository struct {
	DB *pgxpool.Pool
}

func NewDeviceRepository(db *pgxpool.Pool) *DeviceRepository {
	return &DeviceRepository{DB: db}
}

func (r *DeviceRepository) RegisterToken(ctx context.Context, userID, token string) error {
	query := `
		INSERT INTO user_devices (user_id, device_token)
		VALUES ($1, $2)
		ON CONFLICT (user_id, device_token) DO UPDATE 
		SET last_updated = CURRENT_TIMESTAMP
	`
	_, err := r.DB.Exec(ctx, query, userID, token)
	return err
}
