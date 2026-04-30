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

// InventoryClient calls the inventory service on behalf of the booking service.
type InventoryClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewInventoryClient(baseURL string) *InventoryClient {
	return &InventoryClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type seatIDsPayload struct {
	SeatIDs []string `json:"seatIds"`
}

func (c *InventoryClient) ConfirmSeats(ctx context.Context, showtimeID string, seatIDs []string) error {
	body, err := json.Marshal(seatIDsPayload{SeatIDs: seatIDs})
	if err != nil {
		return apperrors.Wrap(apperrors.CodeInternal, "marshal inventory request", err)
	}

	url := fmt.Sprintf("%s/api/inventory/internal/showtimes/%s/seats/confirm", c.baseURL, showtimeID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return apperrors.Wrap(apperrors.CodeInternal, "build inventory request", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return apperrors.Wrap(apperrors.CodeInternal, "call inventory service", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return apperrors.New(apperrors.CodeInternal, fmt.Sprintf("inventory service returned %d for confirm", resp.StatusCode))
	}
	return nil
}

func (c *InventoryClient) ReleaseSeats(ctx context.Context, showtimeID string, seatIDs []string) error {
	body, err := json.Marshal(seatIDsPayload{SeatIDs: seatIDs})
	if err != nil {
		return apperrors.Wrap(apperrors.CodeInternal, "marshal inventory request", err)
	}

	url := fmt.Sprintf("%s/api/inventory/internal/showtimes/%s/seats/release", c.baseURL, showtimeID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return apperrors.Wrap(apperrors.CodeInternal, "build inventory request", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return apperrors.Wrap(apperrors.CodeInternal, "call inventory service", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return apperrors.New(apperrors.CodeInternal, fmt.Sprintf("inventory service returned %d for release", resp.StatusCode))
	}
	return nil
}
