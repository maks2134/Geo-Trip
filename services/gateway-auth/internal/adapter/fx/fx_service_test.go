package fx

import (
	"context"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFXService_GetRates_Fallback(t *testing.T) {
	// Service pointing to non-existent redis to test fallback cleanly
	svc := NewFXService("localhost:59999")

	rates := svc.GetRates(context.Background(), "EUR")
	if rates.Base != "EUR" {
		t.Fatalf("expected base EUR, got %s", rates.Base)
	}
	if rates.Rates["USD"] <= 0 {
		t.Fatalf("expected positive USD rate, got %f", rates.Rates["USD"])
	}
	if rates.Rates["EUR"] != 1.0 {
		t.Fatalf("expected EUR rate = 1.0, got %f", rates.Rates["EUR"])
	}
	if rates.Source != "fallback" {
		t.Fatalf("expected source fallback, got %s", rates.Source)
	}

	// Second query should hit memory cache
	cached := svc.GetRates(context.Background(), "EUR")
	if cached.Source != "cache" {
		t.Fatalf("expected source cache, got %s", cached.Source)
	}
}

func TestFXService_Rebase(t *testing.T) {
	eurRates := map[string]float64{
		"EUR": 1.0,
		"USD": 1.10,
		"JPY": 165.0,
	}

	usdRates := rebaseRates(eurRates, "USD")
	if usdRates["USD"] != 1.0 {
		t.Errorf("expected USD rate 1.0, got %f", usdRates["USD"])
	}

	expectedJPYinUSD := 165.0 / 1.10
	if math.Abs(usdRates["JPY"]-expectedJPYinUSD) > 0.001 {
		t.Errorf("expected JPY rate %f, got %f", expectedJPYinUSD, usdRates["JPY"])
	}

	expectedEURinUSD := 1.0 / 1.10
	if math.Abs(usdRates["EUR"]-expectedEURinUSD) > 0.001 {
		t.Errorf("expected EUR rate %f, got %f", expectedEURinUSD, usdRates["EUR"])
	}
}

func TestFXService_ServeHTTP(t *testing.T) {
	svc := NewFXService("")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/fx/rates?base=USD", nil)
	rec := httptest.NewRecorder()

	svc.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected HTTP 200, got %d", rec.Code)
	}

	var resp FXRatesResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Base != "USD" {
		t.Errorf("expected base USD, got %s", resp.Base)
	}
	if resp.Rates["USD"] != 1.0 {
		t.Errorf("expected USD rate 1.0, got %f", resp.Rates["USD"])
	}
}
