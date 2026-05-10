package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	apperrors "github.com/shownest/pkg/errors"
)

type InventoryClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewInventoryClient(baseURL string) *InventoryClient {
	return &InventoryClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

type SeatPriceInfo struct {
	SeatID          string `json:"seatId"`
	CategoryID      string `json:"categoryId"`
	PriceMultiplier string `json:"priceMultiplier"`
}

type OccupancyInfo struct {
	TotalSeats       int     `json:"totalSeats"`
	BookedSeats      int     `json:"bookedSeats"`
	OccupancyPercent float64 `json:"occupancyPercent"`
}

func (c *InventoryClient) GetSeatPrices(ctx context.Context, showtimeID string, seatIDs []string) ([]SeatPriceInfo, error) {
	body, err := json.Marshal(struct {
		SeatIDs []string `json:"seatIds"`
	}{SeatIDs: seatIDs})
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "marshal seat prices request", err)
	}

	url := fmt.Sprintf("%s/api/inventory/internal/showtimes/%s/seats/prices", c.baseURL, showtimeID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build inventory seat prices request", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "call inventory service for prices", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, apperrors.New(apperrors.CodeInternal, fmt.Sprintf("inventory service returned %d for seat prices", resp.StatusCode))
	}

	var prices []SeatPriceInfo
	if err := json.NewDecoder(resp.Body).Decode(&prices); err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "decode seat prices response", err)
	}
	return prices, nil
}

func (c *InventoryClient) GetOccupancy(ctx context.Context, showtimeID string) (*OccupancyInfo, error) {
	url := fmt.Sprintf("%s/api/inventory/internal/showtimes/%s/occupancy", c.baseURL, showtimeID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build inventory occupancy request", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "call inventory service for occupancy", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, apperrors.New(apperrors.CodeInternal, fmt.Sprintf("inventory service returned %d for occupancy", resp.StatusCode))
	}

	var info OccupancyInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "decode occupancy response", err)
	}
	return &info, nil
}
