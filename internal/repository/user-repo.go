package repository

import (
	"context"
	"errors"

	"github.com/NoahFola/travel_app_backend/internal/domain" // Replace with your actual module name
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{DB: db}
}

// Create inserts a new user and scans the generated UUID back into the struct
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (email, password_hash, auth_provider, created_at, updated_at) 
		VALUES ($1, $2, $3, NOW(), NOW()) 
		RETURNING id, created_at`

	// default provider to email if empty
	if user.AuthProvider == "" {
		user.AuthProvider = "email"
	}

	// We pass pointers for Email and PasswordHash because they can be nil in the struct (though required for this specific query)
	err := r.DB.QueryRow(ctx, query, user.Email, user.PasswordHash, user.AuthProvider).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

// GetByEmail finds a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, avatar_url, auth_provider, created_at 
		FROM users 
		WHERE email = $1`

	var user domain.User
	err := r.DB.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.AvatarURL,
		&user.AuthProvider,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// UpsertOAuthUser creates a user if they don't exist, or updates them if they do
func (r *UserRepository) UpsertOAuthUser(ctx context.Context, email, provider, providerID, name, avatar string) (*domain.User, error) {
	// ON CONFLICT(email): If user exists with this email, just update their info.
	// We trust Google/Facebook emails are verified.
	query := `
		INSERT INTO users (email, auth_provider, provider_user_id, full_name, avatar_url, email_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, true, NOW(), NOW())
		ON CONFLICT (email) 
		DO UPDATE SET 
			full_name = EXCLUDED.full_name,
			avatar_url = EXCLUDED.avatar_url,
			provider_user_id = EXCLUDED.provider_user_id, -- Link the account
			auth_provider = EXCLUDED.auth_provider,
			updated_at = NOW()
		RETURNING id, email, created_at`

	var user domain.User
	// Scan the result
	err := r.DB.QueryRow(ctx, query, email, provider, providerID, name, avatar).Scan(&user.ID, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
