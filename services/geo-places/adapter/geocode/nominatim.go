package geocode

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/makskozlov/atlasflow/services/geo-places/domain"
)

type NominatimGeocoder struct {
	baseURL    string
	httpClient *http.Client
	userAgent  string
}

func NewNominatimGeocoder(baseURL string, client *http.Client) *NominatimGeocoder {
	if baseURL == "" {
		baseURL = "https://nominatim.openstreetmap.org"
	}
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}
	return &NominatimGeocoder{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: client,
		userAgent:  "AtlasFlow/1.0 (contact@atlasflow.internal)",
	}
}

type nominatimResponse struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Address     struct {
		Road        string `json:"road"`
		City        string `json:"city"`
		Town        string `json:"town"`
		Village     string `json:"village"`
		CountryCode string `json:"country_code"`
		Country     string `json:"country"`
	} `json:"address"`
}

func (n *NominatimGeocoder) ReverseGeocode(ctx context.Context, loc domain.LatLng) (*domain.ReverseGeocodeMetadata, error) {
	if err := loc.Validate(); err != nil {
		return nil, err
	}

	u, err := url.Parse(n.baseURL + "/reverse")
	if err != nil {
		return nil, fmt.Errorf("parse nominatim url: %w", err)
	}

	q := u.Query()
	q.Set("format", "jsonv2")
	q.Set("lat", strconv.FormatFloat(loc.Latitude, 'f', 6, 64))
	q.Set("lon", strconv.FormatFloat(loc.Longitude, 'f', 6, 64))
	if locVal, ok := ctx.Value("locale").(string); ok && locVal != "" {
		q.Set("accept-language", locVal)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create geocode req: %w", err)
	}
	req.Header.Set("User-Agent", n.userAgent)

	resp, err := n.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrGeocodingFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status %d", domain.ErrGeocodingFailed, resp.StatusCode)
	}

	var data nominatimResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode nominatim response: %w", err)
	}

	city := data.Address.City
	if city == "" {
		city = data.Address.Town
	}
	if city == "" {
		city = data.Address.Village
	}

	name := data.Name
	if name == "" {
		name = data.Address.Road
	}
	if name == "" {
		name = city
	}

	return &domain.ReverseGeocodeMetadata{
		Name:        name,
		Address:     data.DisplayName,
		City:        city,
		CountryISO2: strings.ToUpper(data.Address.CountryCode),
		CountryName: data.Address.Country,
		Location:    loc,
	}, nil
}