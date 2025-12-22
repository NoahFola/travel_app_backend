package handlers

import (
	"net/http"

	"github.com/NoahFola/travel_app_backend/internal/service" // update module name
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Service *service.AuthService
}

// DTOs (Data Transfer Objects)
type authRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type oauthRequest struct {
	Token string `json:"token" binding:"required"`
}

func (h *AuthHandler) Signup(c *gin.Context) {
	var req authRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	access, refresh, err := h.Service.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User likely already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"access_token": access, "refresh_token": refresh})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req authRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	access, refresh, err := h.Service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": access, "refresh_token": refresh})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token is required"})
		return
	}

	newAccess, err := h.Service.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": newAccess})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Since JWTs are stateless, the server doesn't "delete" anything unless you use a Redis blacklist.
	// For MVP, the Client is responsible for deleting the token.
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	var req oauthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	access, refresh, err := h.Service.LoginWithGoogle(c.Request.Context(), req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Google authentication failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": access, "refresh_token": refresh})
}
