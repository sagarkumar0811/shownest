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

// InventoryClient forwards the user's JWT so inventory can validate ownership of the seat locks.
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

func (c *InventoryClient) ConfirmSeats(ctx context.Context, authHeader, showtimeID string, seatIDs []string) error {
	return c.postSeatAction(ctx, authHeader, showtimeID, "confirm", seatIDs)
}

func (c *InventoryClient) ReleaseSeats(ctx context.Context, authHeader, showtimeID string, seatIDs []string) error {
	return c.postSeatAction(ctx, authHeader, showtimeID, "release", seatIDs)
}

func (c *InventoryClient) postSeatAction(ctx context.Context, authHeader, showtimeID, action string, seatIDs []string) error {
	body, err := json.Marshal(seatIDsPayload{SeatIDs: seatIDs})
	if err != nil {
		return apperrors.Wrap(apperrors.CodeInternal, "marshal inventory request", err)
	}

	url := fmt.Sprintf("%s/api/inventory/v1/showtimes/%s/seats/%s", c.baseURL, showtimeID, action)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return apperrors.Wrap(apperrors.CodeInternal, "build inventory request", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return apperrors.Wrap(apperrors.CodeInternal, "call inventory service", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return apperrors.New(apperrors.CodeInternal, fmt.Sprintf("inventory service returned %d for action %s", resp.StatusCode, action))
	}
	return nil
}
