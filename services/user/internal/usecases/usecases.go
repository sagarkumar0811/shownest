package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	pkgaws "github.com/shownest/pkg/aws"
	apperrors "github.com/shownest/pkg/errors"
	"github.com/shownest/pkg/jwt"
	"github.com/shownest/pkg/logger"
	"github.com/shownest/user-service/internal/dto/request"
	"github.com/shownest/user-service/internal/dto/response"
	"github.com/shownest/user-service/internal/mapper"
	"github.com/shownest/user-service/internal/models"
	"github.com/shownest/user-service/internal/repository"
	"github.com/shownest/user-service/internal/utils"
	"go.uber.org/zap"
)

type UseCase struct {
	repo   *repository.Repository
	cache  *redis.Client
	sns    *pkgaws.SNSClient
	jwtSvc *jwt.Service
}

func New(repo *repository.Repository, cache *redis.Client, sns *pkgaws.SNSClient, jwtSvc *jwt.Service) *UseCase {
	return &UseCase{repo: repo, cache: cache, sns: sns, jwtSvc: jwtSvc}
}

func (uc *UseCase) SendOTP(ctx context.Context, phone string) error {
	if !utils.IsValidPhone(phone) {
		return apperrors.New(apperrors.CodeInvalidArgument, "phone must be in E.164 format (e.g. +919876543210)")
	}

	attKey := utils.OTPAttemptsPrefix + phone
	count, err := uc.cache.Incr(ctx, attKey).Result()
	if err != nil {
		return apperrors.Wrap(apperrors.CodeInternal, "rate limit check failed", err)
	}
	if count == 1 {
		uc.cache.Expire(ctx, attKey, utils.OTPAttemptsTTL)
	}
	if count > utils.MaxOTPAttempts {
		return apperrors.New(apperrors.CodeResourceExhausted, "too many OTP requests; try again later")
	}

	otp, err := utils.GenerateOTP()
	if err != nil {
		return apperrors.Wrap(apperrors.CodeInternal, "generate otp", err)
	}

	hashed := utils.HashSHA256(otp)
	otpKey := utils.OTPKeyPrefix + phone
	if err := uc.cache.Set(ctx, otpKey, hashed, utils.OTPTTL).Err(); err != nil {
		return apperrors.Wrap(apperrors.CodeInternal, "store otp", err)
	}

	message := fmt.Sprintf("Your ShowNest OTP is %s. Valid for 5 minutes. Do not share it.", otp)
	if err := uc.sns.SendSMS(ctx, phone, message); err != nil {
		return apperrors.Wrap(apperrors.CodeInternal, "send otp sms", err)
	}

	logger.WithContext(ctx).Info("OTP sent", zap.String("phone", phone))
	return nil
}

func (uc *UseCase) VerifyOTP(ctx context.Context, req request.VerifyOTPRequest, deviceInfo, ipAddress string) (*response.AuthResponse, error) {
	if !utils.IsValidPhone(req.Phone) {
		return nil, apperrors.New(apperrors.CodeInvalidArgument, "invalid phone format")
	}

	otpKey := utils.OTPKeyPrefix + req.Phone
	storedHash, err := uc.cache.Get(ctx, otpKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, apperrors.New(apperrors.CodeInvalidCredentials, "invalid or expired OTP")
		}
		return nil, apperrors.Wrap(apperrors.CodeInternal, "fetch otp from cache", err)
	}

	if utils.HashSHA256(req.OTP) != storedHash {
		return nil, apperrors.New(apperrors.CodeInvalidCredentials, "invalid or expired OTP")
	}

	uc.cache.Del(ctx, otpKey)

	user, err := uc.repo.GetUserByPhone(ctx, req.Phone)
	if err != nil {
		if apperrors.HasCode(err, apperrors.CodeDBNotFound) {
			user, err = uc.repo.CreateUser(ctx, req.Phone, utils.UserRoleCustomer)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if user.Status == utils.UserStatusBlocked {
		return nil, apperrors.New(apperrors.CodeUserBlocked, "account is blocked")
	}

	return uc.issueTokenPair(ctx, user, deviceInfo, ipAddress)
}

func (uc *UseCase) RefreshToken(ctx context.Context, rawRefreshToken string) (*response.AuthResponse, error) {
	claims, err := uc.jwtSvc.ValidateRefreshToken(rawRefreshToken)
	if err != nil {
		return nil, apperrors.New(apperrors.CodeTokenInvalid, "invalid refresh token")
	}

	tokenHash := utils.HashSHA256(rawRefreshToken)
	session, err := uc.repo.GetSessionByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, apperrors.New(apperrors.CodeTokenInvalid, "session not found or expired")
	}

	user, err := uc.repo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}
	if user.Status == utils.UserStatusBlocked {
		return nil, apperrors.New(apperrors.CodeUserBlocked, "account is blocked")
	}

	newRefreshToken, err := uc.jwtSvc.GenerateRefreshToken(user.ID, user.Role, session.ID)
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "generate refresh token", err)
	}
	newTokenHash := utils.HashSHA256(newRefreshToken)
	newExpiresAt := time.Now().Add(utils.RefreshTokenDuration)

	if err := uc.repo.RotateSessionToken(ctx, session.ID, newTokenHash, newExpiresAt); err != nil {
		return nil, err
	}

	accessToken, err := uc.jwtSvc.GenerateAccessToken(user.ID, user.Role, session.ID)
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "generate access token", err)
	}

	return mapper.ToAuthResponse(accessToken, newRefreshToken, user), nil
}

func (uc *UseCase) ListSessions(ctx context.Context, userID string) ([]response.SessionInfo, error) {
	sessions, err := uc.repo.ListActiveSessions(ctx, userID)
	if err != nil {
		return nil, err
	}
	return mapper.ToSessionInfoList(sessions), nil
}

func (uc *UseCase) RevokeSession(ctx context.Context, userID, sessionID string) error {
	return uc.repo.RevokeSession(ctx, sessionID, userID)
}

func (uc *UseCase) RevokeAllSessions(ctx context.Context, userID, currentSessionID string) error {
	return uc.repo.RevokeAllSessionsExcept(ctx, userID, currentSessionID)
}

func (uc *UseCase) GetProfile(ctx context.Context, userID string) (*response.UserInfo, error) {
	user, err := uc.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	info := mapper.ToUserInfo(user)
	return &info, nil
}

func (uc *UseCase) UpdateProfile(ctx context.Context, userID string, req request.UpdateProfileRequest) (*response.UserInfo, error) {
	if err := uc.repo.UpdateUserEmail(ctx, userID, req.Email); err != nil {
		return nil, err
	}
	return uc.GetProfile(ctx, userID)
}

func (uc *UseCase) issueTokenPair(ctx context.Context, user *models.User, deviceInfo, ipAddress string) (*response.AuthResponse, error) {
	sessionID, err := utils.NewUUID()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "generate session id", err)
	}

	accessToken, err := uc.jwtSvc.GenerateAccessToken(user.ID, user.Role, sessionID)
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "generate access token", err)
	}

	refreshToken, err := uc.jwtSvc.GenerateRefreshToken(user.ID, user.Role, sessionID)
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "generate refresh token", err)
	}

	expiresAt := time.Now().Add(utils.RefreshTokenDuration)
	if _, err := uc.repo.CreateSession(ctx, sessionID, user.ID, utils.HashSHA256(refreshToken), deviceInfo, ipAddress, expiresAt); err != nil {
		return nil, err
	}

	return mapper.ToAuthResponse(accessToken, refreshToken, user), nil
}
