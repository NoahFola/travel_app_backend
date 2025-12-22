package domain

import (
	"time"
)

type User struct {
	ID             string     `json:"id"`    // UUIDs are handled as strings in Go
	Email          *string    `json:"email"` // Pointer because it can be null (OAuth)
	EmailVerified  bool       `json:"email_verified"`
	PasswordHash   *string    `json:"-"` // Pointer (null if OAuth)
	FullName       *string    `json:"full_name"`
	AvatarURL      *string    `json:"avatar_url"`
	AuthProvider   string     `json:"auth_provider"`
	ProviderUserID *string    `json:"provider_user_id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	LastLoginAt    *time.Time `json:"last_login_at"`
}
