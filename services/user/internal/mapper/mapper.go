package mapper

import (
	"github.com/shownest/user-service/internal/dto/response"
	"github.com/shownest/user-service/internal/models"
)

func ToUserInfo(u *models.User) response.UserInfo {
	return response.UserInfo{
		ID:    u.ID,
		Phone: u.Phone,
		Email: u.Email,
		Role:  u.Role,
	}
}

func ToSessionInfoList(sessions []models.Session) []response.SessionInfo {
	infos := make([]response.SessionInfo, len(sessions))
	for i, s := range sessions {
		infos[i] = response.SessionInfo{
			ID:         s.ID,
			DeviceInfo: s.DeviceInfo,
			IPAddress:  s.IPAddress,
			CreatedAt:  s.CreatedAt,
			ExpiresAt:  s.ExpiresAt,
		}
	}
	return infos
}

func ToAuthResponse(accessToken, refreshToken string, user *models.User) *response.AuthResponse {
	return &response.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         ToUserInfo(user),
	}
}
