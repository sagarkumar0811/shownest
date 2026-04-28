package jwt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shownest/pkg/config"
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

// Config holds the configuration for JWT token generation and validation.
type Config struct {
	AccessSecret  string        `json:"accessSecret"`
	RefreshSecret string        `json:"refreshSecret"`
	AccessExpiry  time.Duration `json:"accessExpiry"`
	RefreshExpiry time.Duration `json:"refreshExpiry"`
}

// UnmarshalJSON custom unmarshals the Config from JSON, parsing duration strings into time.Duration.
func (c *Config) UnmarshalJSON(data []byte) error {
	var raw struct {
		AccessSecret  string `json:"accessSecret"`
		RefreshSecret string `json:"refreshSecret"`
		AccessExpiry  string `json:"accessExpiry"`
		RefreshExpiry string `json:"refreshExpiry"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	var err error
	c.AccessSecret = raw.AccessSecret
	c.RefreshSecret = raw.RefreshSecret
	if c.AccessExpiry, err = time.ParseDuration(raw.AccessExpiry); err != nil {
		return fmt.Errorf("parse accessExpiry %q: %w", raw.AccessExpiry, err)
	}
	if c.RefreshExpiry, err = time.ParseDuration(raw.RefreshExpiry); err != nil {
		return fmt.Errorf("parse refreshExpiry %q: %w", raw.RefreshExpiry, err)
	}
	return nil
}

// Init initializes the JWT service by loading the configuration from the provided config provider.
func Init(ctx context.Context, provider config.ConfigProvider) (*Service, error) {
	raw, err := provider.Get(ctx, config.JWTConfig)
	if err != nil {
		return nil, fmt.Errorf("jwt: get config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return nil, fmt.Errorf("jwt: parse config: %w", err)
	}

	return &Service{Config: cfg}, nil
}

// Service provides methods for generating and validating JWT tokens.
type Service struct {
	Config Config
}

// GenerateAccessToken creates a short-lived access token containing user claims.
func (s *Service) GenerateAccessToken(userID, role, sessionID string) (string, error) {
	return s.generateToken(userID, role, sessionID, s.Config.AccessSecret, s.Config.AccessExpiry)
}

// GenerateRefreshToken creates a long-lived refresh token that can be used to obtain new access tokens without re-authenticating the user.
func (s *Service) GenerateRefreshToken(userID, role, sessionID string) (string, error) {
	return s.generateToken(userID, role, sessionID, s.Config.RefreshSecret, s.Config.RefreshExpiry)
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
	return s.validateToken(tokenString, s.Config.AccessSecret)
}

// ValidateRefreshToken validates a refresh token and returns its claims.
func (s *Service) ValidateRefreshToken(tokenString string) (*Claims, error) {
	return s.validateToken(tokenString, s.Config.RefreshSecret)
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
