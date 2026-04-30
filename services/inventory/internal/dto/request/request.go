package request

import "errors"

type HallIDRequest struct {
	HallID string `uri:"id" binding:"required"`
}

type CreateSeatCategoryRequest struct {
	Name            string  `json:"name"            binding:"required,min=1,max=50"`
	PriceMultiplier float64 `json:"priceMultiplier" binding:"required,gt=0"`
}

type CreateSeatRequest struct {
	CategoryID string   `json:"categoryId" binding:"required"`
	Row        string   `json:"row"        binding:"required,min=1,max=10"`
	Number     int      `json:"number"     binding:"required,min=1"`
	XPosition  *float64 `json:"xPosition"`
	YPosition  *float64 `json:"yPosition"`
}

type BulkCreateSeatsRequest struct {
	Seats []CreateSeatRequest `json:"seats" binding:"required,min=1,dive"`
}

func (r *BulkCreateSeatsRequest) Validate() error {
	if len(r.Seats) == 0 {
		return errors.New("seats must not be empty")
	}
	return nil
}

type ShowtimeIDRequest struct {
	ShowtimeID string `uri:"id" binding:"required"`
}

type PublishShowtimeSeatsRequest struct {
	ShowtimeID string `json:"showtimeId" binding:"required"`
}

type LockSeatsRequest struct {
	SeatIDs []string `json:"seatIds" binding:"required,min=1"`
}

func (r *LockSeatsRequest) Validate() error {
	if len(r.SeatIDs) == 0 {
		return errors.New("seatIds must not be empty")
	}
	return nil
}

type ReleaseSeatsRequest struct {
	SeatIDs []string `json:"seatIds" binding:"required,min=1"`
}

func (r *ReleaseSeatsRequest) Validate() error {
	if len(r.SeatIDs) == 0 {
		return errors.New("seatIds must not be empty")
	}
	return nil
}

type ConfirmSeatsRequest struct {
	SeatIDs []string `json:"seatIds" binding:"required,min=1"`
}

type GetSeatPricesRequest struct {
	SeatIDs []string `json:"seatIds" binding:"required,min=1"`
}
