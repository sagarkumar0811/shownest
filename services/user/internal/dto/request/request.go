package request

import (
	"github.com/shownest/user-service/internal/utils"
)

type SendOTPRequest struct {
	Phone string `json:"phone" binding:"required"`
}

func (r *SendOTPRequest) Validate() error {
	phone, err := utils.FilterMobileNumber(r.Phone)
	if err != nil {
		return err
	}
	r.Phone = phone
	return nil
}

type VerifyOTPRequest struct {
	Phone string `json:"phone" binding:"required"`
	OTP   string `json:"otp"   binding:"required,len=6"`
}

func (r *VerifyOTPRequest) Validate() error {
	phone, err := utils.FilterMobileNumber(r.Phone)
	if err != nil {
		return err
	}
	r.Phone = phone
	return nil
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type RevokeSessionRequest struct {
	SessionID string `uri:"id" binding:"required"`
}

type UpdateProfileRequest struct {
	Email *string `json:"email" binding:"omitempty,email"`
}
