package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	apperrors "github.com/shownest/pkg/errors"
)

type CatalogClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewCatalogClient(baseURL string) *CatalogClient {
	return &CatalogClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

type ShowtimeInfo struct {
	ID        string    `json:"id"`
	EventID   string    `json:"eventId"`
	HallID    string    `json:"hallId"`
	StartTime time.Time `json:"startTime"`
	BasePrice string    `json:"basePrice"`
	Status    string    `json:"status"`
}

type EventInfo struct {
	ID       string `json:"id"`
	Category string `json:"category"`
}

func (c *CatalogClient) GetShowtime(ctx context.Context, showtimeID string) (*ShowtimeInfo, error) {
	url := fmt.Sprintf("%s/api/catalog/v1/showtimes/%s", c.baseURL, showtimeID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build catalog showtime request", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "call catalog service", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, apperrors.New(apperrors.CodeDBNotFound, "showtime not found")
	}
	if resp.StatusCode >= 400 {
		return nil, apperrors.New(apperrors.CodeInternal, fmt.Sprintf("catalog service returned %d for showtime", resp.StatusCode))
	}

	var info ShowtimeInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "decode catalog showtime response", err)
	}
	return &info, nil
}

func (c *CatalogClient) GetEvent(ctx context.Context, eventID string) (*EventInfo, error) {
	url := fmt.Sprintf("%s/api/catalog/v1/events/%s", c.baseURL, eventID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build catalog event request", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "call catalog service", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, apperrors.New(apperrors.CodeDBNotFound, "event not found")
	}
	if resp.StatusCode >= 400 {
		return nil, apperrors.New(apperrors.CodeInternal, fmt.Sprintf("catalog service returned %d for event", resp.StatusCode))
	}

	var info EventInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "decode catalog event response", err)
	}
	return &info, nil
}
