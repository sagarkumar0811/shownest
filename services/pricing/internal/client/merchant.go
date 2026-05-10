package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	apperrors "github.com/shownest/pkg/errors"
)

type MerchantClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewMerchantClient(baseURL string) *MerchantClient {
	return &MerchantClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *MerchantClient) GetMerchantIDByUserID(ctx context.Context, userID string) (string, error) {
	url := fmt.Sprintf("%s/api/merchant/internal/merchants/%s", c.baseURL, userID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", apperrors.Wrap(apperrors.CodeInternal, "build merchant client request", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", apperrors.Wrap(apperrors.CodeInternal, "call merchant service", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", apperrors.New(apperrors.CodeDBNotFound, "merchant not found")
	}
	if resp.StatusCode >= 400 {
		return "", apperrors.New(apperrors.CodeInternal, fmt.Sprintf("merchant service returned %d", resp.StatusCode))
	}

	var body struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", apperrors.Wrap(apperrors.CodeInternal, "decode merchant response", err)
	}

	return body.ID, nil
}
