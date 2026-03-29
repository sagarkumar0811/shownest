package usecases

import (
	"context"
	"time"

	"github.com/shownest/merchant-service/internal/config"
	"github.com/shownest/merchant-service/internal/dto/request"
	"github.com/shownest/merchant-service/internal/dto/response"
	"github.com/shownest/merchant-service/internal/mapper"
	"github.com/shownest/merchant-service/internal/repository"
	"github.com/shownest/merchant-service/internal/utils"
	pkgaws "github.com/shownest/pkg/aws"
	apperrors "github.com/shownest/pkg/errors"
)

type UseCase struct {
	repo   *repository.Repository
	s3     *pkgaws.S3Client
	config *config.Config
}

func New(repo *repository.Repository, s3 *pkgaws.S3Client, cfg *config.Config) *UseCase {
	return &UseCase{repo: repo, s3: s3, config: cfg}
}

func (uc *UseCase) CreateMerchant(ctx context.Context, userID string, req request.CreateMerchantRequest) (*response.MerchantInfo, error) {
	_, err := uc.repo.GetMerchantByUserID(ctx, userID)
	if err == nil {
		return nil, apperrors.New(apperrors.CodeAlreadyExists, "merchant profile already exists")
	}
	if !apperrors.HasCode(err, apperrors.CodeDBNotFound) {
		return nil, err
	}

	merchant, err := uc.repo.CreateMerchant(ctx, userID, req.BusinessName, req.Category, req.ContactPhone, req.ContactEmail)
	if err != nil {
		return nil, err
	}
	info := mapper.ToMerchantInfo(merchant)
	return &info, nil
}

func (uc *UseCase) GetMyMerchant(ctx context.Context, userID string) (*response.MerchantInfo, error) {
	merchant, err := uc.repo.GetMerchantByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	info := mapper.ToMerchantInfo(merchant)
	return &info, nil
}

func (uc *UseCase) SubmitForReview(ctx context.Context, userID string) error {
	merchant, err := uc.repo.GetMerchantByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if merchant.Status != utils.MerchantStatusDraft {
		return apperrors.New(apperrors.CodeFailedPrecondition, "only a draft merchant can be submitted for review")
	}
	return uc.repo.UpdateMerchantStatus(ctx, merchant.ID, utils.MerchantStatusPending)
}

func (uc *UseCase) CreateVenue(ctx context.Context, userID string, req request.CreateVenueRequest) (*response.VenueInfo, error) {
	merchant, err := uc.repo.GetMerchantByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	v, err := uc.repo.CreateVenue(ctx, merchant.ID, req.Name, req.Address, req.City, req.State, req.Pincode, req.Latitude, req.Longitude)
	if err != nil {
		return nil, err
	}
	info := mapper.ToVenueInfo(v)
	return &info, nil
}

func (uc *UseCase) ListMyVenues(ctx context.Context, userID string) ([]response.VenueInfo, error) {
	merchant, err := uc.repo.GetMerchantByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	venues, err := uc.repo.GetVenuesByMerchantID(ctx, merchant.ID)
	if err != nil {
		return nil, err
	}
	return mapper.ToVenueInfoList(venues), nil
}

func (uc *UseCase) GetVenue(ctx context.Context, venueID, userID string) (*response.VenueInfo, error) {
	merchant, err := uc.repo.GetMerchantByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	v, err := uc.repo.GetVenueByID(ctx, venueID)
	if err != nil {
		return nil, err
	}
	if v.MerchantID != merchant.ID {
		return nil, apperrors.New(apperrors.CodePermissionDenied, "venue does not belong to your merchant account")
	}
	info := mapper.ToVenueInfo(v)
	return &info, nil
}

func (uc *UseCase) GetNearbyVenues(ctx context.Context, lat, lng, radiusMeters float64) ([]response.VenueWithDistanceInfo, error) {
	venues, err := uc.repo.GetNearbyVenues(ctx, lat, lng, radiusMeters)
	if err != nil {
		return nil, err
	}
	return mapper.ToVenueWithDistanceInfoList(venues), nil
}

func (uc *UseCase) GetVenuesByCity(ctx context.Context, city string) ([]response.VenueInfo, error) {
	venues, err := uc.repo.GetVenuesByCity(ctx, city)
	if err != nil {
		return nil, err
	}
	return mapper.ToVenueInfoList(venues), nil
}

func (uc *UseCase) CreateHall(ctx context.Context, venueID, userID string, req request.CreateHallRequest) (*response.HallInfo, error) {
	merchant, err := uc.repo.GetMerchantByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	v, err := uc.repo.GetVenueByID(ctx, venueID)
	if err != nil {
		return nil, err
	}
	if v.MerchantID != merchant.ID {
		return nil, apperrors.New(apperrors.CodePermissionDenied, "venue does not belong to your merchant account")
	}
	h, err := uc.repo.CreateHall(ctx, venueID, req.Name, req.Capacity, req.HallType)
	if err != nil {
		return nil, err
	}
	info := mapper.ToHallInfo(h)
	return &info, nil
}

func (uc *UseCase) ListHalls(ctx context.Context, venueID, userID string) ([]response.HallInfo, error) {
	merchant, err := uc.repo.GetMerchantByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	v, err := uc.repo.GetVenueByID(ctx, venueID)
	if err != nil {
		return nil, err
	}
	if v.MerchantID != merchant.ID {
		return nil, apperrors.New(apperrors.CodePermissionDenied, "venue does not belong to your merchant account")
	}
	halls, err := uc.repo.GetHallsByVenueID(ctx, venueID)
	if err != nil {
		return nil, err
	}
	return mapper.ToHallInfoList(halls), nil
}

func (uc *UseCase) RequestDocumentUploadURL(ctx context.Context, userID, docType string) (*response.UploadURLResponse, error) {
	merchant, err := uc.repo.GetMerchantByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	id, err := utils.NewUUID()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "generate key id", err)
	}

	s3Key := utils.GetS3Key(uc.config.App, merchant.ID, "documents", docType, id)
	ttl := time.Duration(utils.DocumentUploadURLTTL) * time.Minute
	uploadURL, err := uc.s3.PresignPutURL(ctx, s3Key, ttl)
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "generate upload url", err)
	}
	return &response.UploadURLResponse{
		UploadURL:    uploadURL,
		S3Key:        s3Key,
		DocumentType: docType,
	}, nil
}

func (uc *UseCase) ConfirmDocument(ctx context.Context, userID string, req request.ConfirmDocumentRequest) (*response.DocumentInfo, error) {
	merchant, err := uc.repo.GetMerchantByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	d, err := uc.repo.CreateDocument(ctx, merchant.ID, req.DocumentType, req.S3Key)
	if err != nil {
		return nil, err
	}
	info := mapper.ToDocumentInfo(d)
	return &info, nil
}

func (uc *UseCase) ListDocuments(ctx context.Context, userID string) ([]response.DocumentInfo, error) {
	merchant, err := uc.repo.GetMerchantByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	docs, err := uc.repo.GetDocumentsByMerchantID(ctx, merchant.ID)
	if err != nil {
		return nil, err
	}
	return mapper.ToDocumentInfoList(docs), nil
}
