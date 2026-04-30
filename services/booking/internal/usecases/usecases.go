package usecases

import (
	"context"
	"time"

	"github.com/shownest/booking-service/internal/client"
	"github.com/shownest/booking-service/internal/dto/request"
	"github.com/shownest/booking-service/internal/dto/response"
	"github.com/shownest/booking-service/internal/mapper"
	"github.com/shownest/booking-service/internal/repository"
	"github.com/shownest/booking-service/internal/utils"
	apperrors "github.com/shownest/pkg/errors"
	pkgutils "github.com/shownest/pkg/utils"
)

type UseCase struct {
	repo      *repository.Repository
	inventory *client.InventoryClient
}

func New(repo *repository.Repository, inventory *client.InventoryClient) *UseCase {
	return &UseCase{repo: repo, inventory: inventory}
}

func (uc *UseCase) CreateBooking(ctx context.Context, userID string, req request.CreateBookingRequest) (*response.BookingInfo, error) {
	var totalAmount float64
	items := make([]struct {
		SeatID, CategoryID string
		Price              float64
	}, len(req.Seats))

	for i, s := range req.Seats {
		items[i] = struct {
			SeatID, CategoryID string
			Price              float64
		}{s.SeatID, s.CategoryID, s.Price}
		totalAmount += s.Price
	}

	expiresAt := time.Now().Add(time.Duration(utils.BookingExpirySeconds) * time.Second)

	booking, bookingItems, err := uc.repo.CreateBookingWithItems(ctx, userID, req.ShowtimeID, totalAmount, expiresAt, items)
	if err != nil {
		return nil, err
	}

	info := mapper.ToBookingInfo(booking, bookingItems)
	return &info, nil
}

// ConfirmBooking transitions a pending booking to confirmed. If the inventory
// call fails, the booking remains pending so the user can retry.
func (uc *UseCase) ConfirmBooking(ctx context.Context, bookingID, userID string) (*response.BookingInfo, error) {
	booking, err := uc.repo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}
	if booking.UserID != userID {
		return nil, apperrors.New(apperrors.CodePermissionDenied, "booking does not belong to user")
	}
	if booking.Status != utils.BookingStatusPending {
		return nil, apperrors.New(apperrors.CodeFailedPrecondition, "only pending bookings can be confirmed")
	}
	if time.Now().After(booking.ExpiresAt) {
		return nil, apperrors.New(apperrors.CodeFailedPrecondition, "booking has expired")
	}

	items, err := uc.repo.GetBookingItems(ctx, bookingID)
	if err != nil {
		return nil, err
	}

	seatIDs := make([]string, len(items))
	for i, item := range items {
		seatIDs[i] = item.SeatID
	}

	// Call inventory to atomically move seats from locked → booked.
	if err := uc.inventory.ConfirmSeats(ctx, booking.ShowtimeID, seatIDs); err != nil {
		return nil, err
	}

	qrToken, err := pkgutils.NewHex(16)
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "generate qr token", err)
	}

	updated, err := uc.repo.UpdateBookingStatus(ctx, bookingID, utils.BookingStatusConfirmed, &qrToken)
	if err != nil {
		return nil, err
	}

	info := mapper.ToBookingInfo(updated, items)
	return &info, nil
}

// CancelBooking cancels a pending booking and releases seat locks back to available.
// Confirmed bookings require a refund flow (future phase).
func (uc *UseCase) CancelBooking(ctx context.Context, bookingID, userID string) error {
	booking, err := uc.repo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return err
	}
	if booking.UserID != userID {
		return apperrors.New(apperrors.CodePermissionDenied, "booking does not belong to user")
	}
	if booking.Status != utils.BookingStatusPending {
		return apperrors.New(apperrors.CodeFailedPrecondition, "only pending bookings can be cancelled here")
	}

	items, err := uc.repo.GetBookingItems(ctx, bookingID)
	if err != nil {
		return err
	}

	seatIDs := make([]string, len(items))
	for i, item := range items {
		seatIDs[i] = item.SeatID
	}

	// Release seat locks in inventory so other users can book them.
	if err := uc.inventory.ReleaseSeats(ctx, booking.ShowtimeID, seatIDs); err != nil {
		return err
	}

	_, err = uc.repo.UpdateBookingStatus(ctx, bookingID, utils.BookingStatusCancelled, nil)
	return err
}

func (uc *UseCase) GetBooking(ctx context.Context, bookingID, userID string) (*response.BookingInfo, error) {
	booking, err := uc.repo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}
	if booking.UserID != userID {
		return nil, apperrors.New(apperrors.CodePermissionDenied, "booking does not belong to user")
	}

	items, err := uc.repo.GetBookingItems(ctx, bookingID)
	if err != nil {
		return nil, err
	}

	info := mapper.ToBookingInfo(booking, items)
	return &info, nil
}

func (uc *UseCase) ListUserBookings(ctx context.Context, userID string) ([]response.BookingInfo, error) {
	bookings, err := uc.repo.ListBookingsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return mapper.ToBookingInfoList(bookings), nil
}
