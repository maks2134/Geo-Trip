package fx

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// FXRatesResponse represents exchange rates against a base currency.
type FXRatesResponse struct {
	Base   string             `json:"base"`
	Date   string             `json:"date"`
	Rates  map[string]float64 `json:"rates"`
	Source string             `json:"source"` // "live", "cache", or "fallback"
}

// Fallback EUR-based rates for offline or resilient operation.
var defaultRatesEUR = map[string]float64{
	"EUR": 1.0,
	"USD": 1.08,
	"RUB": 102.5,
	"GBP": 0.85,
	"JPY": 164.2,
	"THB": 39.5,
	"TRY": 37.8,
	"GEL": 2.95,
	"AED": 3.97,
	"CHF": 0.96,
	"CNY": 7.82,
	"CAD": 1.48,
	"AUD": 1.66,
	"SGD": 1.45,
	"IDR": 17200.0,
	"VND": 27400.0,
	"KRW": 1490.0,
	"INR": 91.5,
	"BRL": 6.25,
	"MXN": 21.8,
	"ILS": 4.05,
	"NZD": 1.82,
	"SEK": 11.45,
	"NOK": 11.75,
	"DKK": 7.46,
	"PLN": 4.31,
	"CZK": 25.3,
	"HUF": 396.0,
	"ZAR": 19.8,
}

// FXService fetches and caches foreign exchange rates using Redis + in-memory fallback.
type FXService struct {
	redisAddr  string
	httpClient *http.Client
	mu         sync.RWMutex
	memCache   map[string]cacheItem
}

type cacheItem struct {
	response  FXRatesResponse
	expiresAt time.Time
}

func NewFXService(redisAddr string) *FXService {
	return &FXService{
		redisAddr: redisAddr,
		httpClient: &http.Client{
			Timeout: 3 * time.Second,
		},
		memCache: make(map[string]cacheItem),
	}
}

// GetRates returns rates relative to baseCurrency (e.g. "EUR" or "USD").
func (s *FXService) GetRates(ctx context.Context, base string) FXRatesResponse {
	base = strings.ToUpper(strings.TrimSpace(base))
	if base == "" {
		base = "EUR"
	}

	// 1. Check in-memory cache
	s.mu.RLock()
	item, ok := s.memCache[base]
	s.mu.RUnlock()
	if ok && time.Now().Before(item.expiresAt) {
		res := item.response
		res.Source = "cache"
		return res
	}

	// 2. Check Redis cache if configured
	if s.redisAddr != "" {
		if cachedJSON, err := s.redisGet(ctx, fmt.Sprintf("atlasflow:fx:%s", base)); err == nil && cachedJSON != "" {
			var res FXRatesResponse
			if err := json.Unmarshal([]byte(cachedJSON), &res); err == nil {
				res.Source = "cache"
				s.mu.Lock()
				s.memCache[base] = cacheItem{
					response:  res,
					expiresAt: time.Now().Add(12 * time.Hour),
				}
				s.mu.Unlock()
				return res
			}
		}
	}

	// 3. Attempt external fetch (Frankfurter API with EUR base)
	liveRates, date, err := s.fetchExternalRates(ctx, "EUR")
	if err == nil && len(liveRates) > 0 {
		rebased := rebaseRates(liveRates, base)
		res := FXRatesResponse{
			Base:   base,
			Date:   date,
			Rates:  rebased,
			Source: "live",
		}

		// Cache in memory
		s.mu.Lock()
		s.memCache[base] = cacheItem{
			response:  res,
			expiresAt: time.Now().Add(12 * time.Hour),
		}
		s.mu.Unlock()

		// Cache in Redis (12h TTL = 43200s)
		if s.redisAddr != "" {
			if data, err := json.Marshal(res); err == nil {
				_ = s.redisSetEX(ctx, fmt.Sprintf("atlasflow:fx:%s", base), string(data), 43200)
			}
		}

		return res
	}

	// 4. Fallback to embedded default rates
	rebased := rebaseRates(defaultRatesEUR, base)
	res := FXRatesResponse{
		Base:   base,
		Date:   time.Now().Format("2006-01-02"),
		Rates:  rebased,
		Source: "fallback",
	}

	s.mu.Lock()
	s.memCache[base] = cacheItem{
		response:  res,
		expiresAt: time.Now().Add(1 * time.Hour),
	}
	s.mu.Unlock()

	return res
}

func (s *FXService) fetchExternalRates(ctx context.Context, base string) (map[string]float64, string, error) {
	url := fmt.Sprintf("https://api.frankfurter.app/latest?base=%s", base)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, "", err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var data struct {
		Amount float64            `json:"amount"`
		Base   string             `json:"base"`
		Date   string             `json:"date"`
		Rates  map[string]float64 `json:"rates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, "", err
	}

	rates := data.Rates
	if rates == nil {
		rates = make(map[string]float64)
	}
	rates[base] = 1.0

	return rates, data.Date, nil
}

func rebaseRates(eurRates map[string]float64, targetBase string) map[string]float64 {
	targetBase = strings.ToUpper(targetBase)
	baseRate, ok := eurRates[targetBase]
	if !ok || baseRate <= 0 {
		baseRate = 1.0
	}

	result := make(map[string]float64, len(eurRates))
	for cur, rate := range eurRates {
		// rate is X per EUR
		// targetBase is Y per EUR
		// X per 1 target = (X / Y)
		result[cur] = rate / baseRate
	}
	result[targetBase] = 1.0
	return result
}

// Minimal RESP protocol implementation for Redis GET
func (s *FXService) redisGet(ctx context.Context, key string) (string, error) {
	d := net.Dialer{Timeout: 500 * time.Millisecond}
	conn, err := d.DialContext(ctx, "tcp", s.redisAddr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(1 * time.Second))
	cmd := fmt.Sprintf("*2\r\n$3\r\nGET\r\n$%d\r\n%s\r\n", len(key), key)
	if _, err := conn.Write([]byte(cmd)); err != nil {
		return "", err
	}

	r := bufio.NewReader(conn)
	prefix, err := r.ReadByte()
	if err != nil {
		return "", err
	}

	if prefix == '$' {
		line, err := r.ReadString('\n')
		if err != nil {
			return "", err
		}
		var length int
		if _, err := fmt.Sscanf(strings.TrimSpace(line), "%d", &length); err != nil {
			return "", err
		}
		if length == -1 {
			return "", nil // Nil bulk string (key does not exist)
		}
		buf := make([]byte, length+2)
		if _, err := io.ReadFull(r, buf); err != nil {
			return "", err
		}
		return string(buf[:length]), nil
	}

	return "", fmt.Errorf("unexpected redis response prefix: %c", prefix)
}

// Minimal RESP protocol implementation for Redis SET with EX seconds
func (s *FXService) redisSetEX(ctx context.Context, key, val string, seconds int) error {
	d := net.Dialer{Timeout: 500 * time.Millisecond}
	conn, err := d.DialContext(ctx, "tcp", s.redisAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(1 * time.Second))
	secStr := fmt.Sprintf("%d", seconds)
	cmd := fmt.Sprintf("*5\r\n$3\r\nSET\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n$2\r\nEX\r\n$%d\r\n%s\r\n",
		len(key), key, len(val), val, len(secStr), secStr)
	if _, err := conn.Write([]byte(cmd)); err != nil {
		return err
	}

	r := bufio.NewReader(conn)
	line, err := r.ReadString('\n')
	if err != nil {
		return err
	}
	if !strings.HasPrefix(line, "+OK") {
		return fmt.Errorf("redis error response: %s", line)
	}
	return nil
}

// ServeHTTP handles GET /api/v1/fx/rates?base=USD
func (s *FXService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	base := r.URL.Query().Get("base")
	if base == "" {
		base = "EUR"
	}

	rates := s.GetRates(r.Context(), base)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	_ = json.NewEncoder(w).Encode(rates)
}
