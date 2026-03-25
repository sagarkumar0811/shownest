package request

type SendOTPRequest struct {
	Phone string `json:"phone" binding:"required,len=10"`
}

type VerifyOTPRequest struct {
	Phone string `json:"phone" binding:"required,len=10"`
	OTP   string `json:"otp"   binding:"required,len=6"`
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
