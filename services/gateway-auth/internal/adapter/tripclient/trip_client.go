package tripclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/domain"
)

type HTTPTripClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewHTTPTripClient(baseURL string) *HTTPTripClient {
	return &HTTPTripClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type tripCoreShareResponse struct {
	Trip *domain.PublicTrip `json:"trip"`
}

func (c *HTTPTripClient) GetPublicTrip(ctx context.Context, slug string) (*domain.PublicTrip, error) {
	url := fmt.Sprintf("%s/api/v1/share/%s", c.baseURL, slug)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create share request: %w", err)
	}

	// Internal gateway client never sets client-spoofed user id on public endpoints
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute share request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, domain.ErrTripNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("trip-core returned HTTP %d", resp.StatusCode)
	}

	var res tripCoreShareResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decode share response: %w", err)
	}

	if res.Trip == nil {
		return nil, domain.ErrTripNotFound
	}
	return res.Trip, nil
}

func (c *HTTPTripClient) ClonePublicTrip(ctx context.Context, slug, targetOwnerID string) (*domain.ClonedTripSummary, error) {
	url := fmt.Sprintf("%s/api/v1/share/%s/clone", c.baseURL, slug)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create clone request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-User-Id", targetOwnerID)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute clone request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, domain.ErrTripNotFound
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("trip-core returned HTTP %d", resp.StatusCode)
	}

	var res struct {
		Trip *domain.ClonedTripSummary `json:"trip"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decode clone response: %w", err)
	}
	if res.Trip == nil {
		return nil, domain.ErrTripNotFound
	}
	return res.Trip, nil
}
