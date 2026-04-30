package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	apperrors "github.com/shownest/pkg/errors"
)

// CatalogClient calls the catalog service on behalf of the booking service.
type CatalogClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewCatalogClient(baseURL string) *CatalogClient {
	return &CatalogClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetShowtimeBasePrice fetches the base price for a showtime from the catalog service.
func (c *CatalogClient) GetShowtimeBasePrice(ctx context.Context, showtimeID string) (float64, error) {
	url := fmt.Sprintf("%s/api/catalog/internal/showtimes/%s/base-price", c.baseURL, showtimeID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, apperrors.Wrap(apperrors.CodeInternal, "build catalog request", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, apperrors.Wrap(apperrors.CodeInternal, "call catalog service", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return 0, apperrors.New(apperrors.CodeInternal, fmt.Sprintf("catalog service returned %d for base price", resp.StatusCode))
	}

	var result struct {
		BasePrice float64 `json:"basePrice"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, apperrors.Wrap(apperrors.CodeInternal, "decode catalog base price response", err)
	}
	return result.BasePrice, nil
}
