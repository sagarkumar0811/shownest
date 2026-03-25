package response

import "time"

type AuthResponse struct {
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
	User         UserInfo `json:"user"`
}

type UserInfo struct {
	ID    string  `json:"id"`
	Phone string  `json:"phone"`
	Email *string `json:"email,omitempty"`
	Role  string  `json:"role"`
}

type SessionInfo struct {
	ID         string    `json:"id"`
	DeviceInfo *string   `json:"deviceInfo,omitempty"`
	IPAddress  *string   `json:"ipAddress,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	ExpiresAt  time.Time `json:"expiresAt"`
}
