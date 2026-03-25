package utils

import "time"

const (
	OTPKeyPrefix      = "otp:"
	OTPAttemptsPrefix = "otp:attempts:"
	OTPTTL            = 5 * time.Minute
	OTPAttemptsTTL    = 1 * time.Hour
	MaxOTPAttempts    = 5

	RefreshTokenDuration = 30 * 24 * time.Hour
)

const (
	UserStatusActive  = "Active"
	UserStatusBlocked = "Blocked"
)

const (
	UserRoleCustomer = "Customer"
	UserRoleMerchant = "Merchant"
)
