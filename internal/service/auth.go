package service

import (
	"context"
	"errors"

	"google.golang.org/api/idtoken"

	"github.com/NoahFola/travel_app_backend/internal/domain"     // update module name
	"github.com/NoahFola/travel_app_backend/internal/repository" // update module name
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Repo *repository.UserRepository
}

func (s *AuthService) Register(ctx context.Context, email, password, name string) (string, string, domain.User, error) {
	// 1. Hash Password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", domain.User{}, err
	}
	hashStr := string(hashed)

	// 2. Create User
	user := &domain.User{
		Email:        &email,
		PasswordHash: &hashStr,
		FullName:     &name,
	}

	if err := s.Repo.Create(ctx, user); err != nil {
		return "", "", domain.User{}, err // Likely duplicate email
	}

	// 3. Generate Tokens
	accessToken, RefreshToken, err := GenerateTokens(user.ID)
	if err != nil {
		return "", "", domain.User{}, err
	}
	return accessToken, RefreshToken, *user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, domain.User, error) {
	// 1. Find User
	user, err := s.Repo.GetByEmail(ctx, email)
	if err != nil {
		return "", "", domain.User{}, errors.New("invalid credentials")
	}

	// 2. Check Password
	if user.PasswordHash == nil {
		return "", "", domain.User{}, errors.New("user uses OAuth")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password)); err != nil {
		return "", "", domain.User{}, errors.New("invalid credentials")
	}

	// 3. Generate Tokens

	accessToken, RefreshToken, err := GenerateTokens(user.ID)
	if err != nil {
		return "", "", domain.User{}, err
	}
	return accessToken, RefreshToken, *user, nil
}

func (s *AuthService) RefreshToken(oldRefreshToken string) (string, error) {
	// 1. Validate the old refresh token
	claims, err := ValidateToken(oldRefreshToken)
	if err != nil {
		return "", err
	}

	if claims.Type != RefreshToken {
		return "", errors.New("invalid token type")
	}

	// 2. Generate NEW Access Token (15 mins)
	// Note: In a stricter system, we would check if the user is banned in DB here
	newAccess, _, err := GenerateTokens(claims.UserID)
	return newAccess, err
}

func (s *AuthService) LoginWithGoogle(ctx context.Context, googleIDToken string) (string, string, error) {
	// 1. Verify the token with Google
	// ClientID is optional here if you want to skip audience check, but recommended for security
	payload, err := idtoken.Validate(ctx, googleIDToken, "")
	if err != nil {
		return "", "", errors.New("invalid google token")
	}

	// 2. Extract Info
	email := payload.Claims["email"].(string)
	name := payload.Claims["name"].(string)
	picture := payload.Claims["picture"].(string)
	sub := payload.Subject // The unique Google User ID

	// 3. Find or Create User in DB
	user, err := s.Repo.UpsertOAuthUser(ctx, email, "google", sub, name, picture)
	if err != nil {
		return "", "", err
	}

	// 4. Issue OUR App's JWTs
	return GenerateTokens(user.ID)
}
