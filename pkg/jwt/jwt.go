package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Errors related to JWT token validation.
var (
	ErrInvalidToken         = errors.New("invalid token")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrTokenExpired         = errors.New("token has expired")
	ErrInvalidClaims        = errors.New("invalid token claims")
)

// Claims represents the JWT claims containing user information.
type Claims struct {
	UserID    string `json:"userId"`
	Role      string `json:"role"`
	SessionID string `json:"sessionId"`
	jwt.RegisteredClaims
}

// Service handles JWT token operations.
type Service struct {
	accessSecret         string
	refreshSecret        string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

// NewService creates a new JWT service with the provided secrets and token durations.
func NewService(accessSecret, refreshSecret string, accessTokenDuration, refreshTokenDuration time.Duration) *Service {
	return &Service{
		accessSecret:         accessSecret,
		refreshSecret:        refreshSecret,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

// GenerateAccessToken creates a short-lived access token containing user claims.
func (s *Service) GenerateAccessToken(userID, role, sessionID string) (string, error) {
	return s.generateToken(userID, role, sessionID, s.accessSecret, s.accessTokenDuration)
}

// GenerateRefreshToken creates a long-lived refresh token that can be used to obtain new access tokens without re-authenticating the user.
func (s *Service) GenerateRefreshToken(userID, role, sessionID string) (string, error) {
	return s.generateToken(userID, role, sessionID, s.refreshSecret, s.refreshTokenDuration)
}

// generateToken is a helper function to create a JWT token with the specified claims, secret, and duration.
func (s *Service) generateToken(userID, role, sessionID, secret string, duration time.Duration) (string, error) {
	claims := &Claims{
		UserID:    userID,
		Role:      role,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateAccessToken validates an access token and returns its claims.
func (s *Service) ValidateAccessToken(tokenString string) (*Claims, error) {
	return s.validateToken(tokenString, s.accessSecret)
}

// ValidateRefreshToken validates a refresh token and returns its claims.
func (s *Service) ValidateRefreshToken(tokenString string) (*Claims, error) {
	return s.validateToken(tokenString, s.refreshSecret)
}

// validateToken is a helper function to validate a token with the given secret and return its claims.
func (s *Service) validateToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidClaims
	}

	return claims, nil
}
