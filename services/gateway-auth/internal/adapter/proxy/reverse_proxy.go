package proxy

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/adapter/middleware"
)

type ReverseProxyHandler struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
}

func NewReverseProxy(targetURL string, logger *slog.Logger) (*ReverseProxyHandler, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	p := httputil.NewSingleHostReverseProxy(target)
	defaultDirector := p.Director

	p.Director = func(req *http.Request) {
		defaultDirector(req)
		req.Host = target.Host

		// Strip any spoofed headers incoming from external clients
		req.Header.Del(middleware.HeaderUserID)
		req.Header.Del(middleware.HeaderUserEmail)

		// Forward verified identity from authenticated context
		if identity, ok := middleware.UserFromContext(req.Context()); ok {
			req.Header.Set(middleware.HeaderUserID, identity.ID)
			req.Header.Set(middleware.HeaderUserEmail, identity.Email)
		}

		// Forward client locale preference (MH-A)
		if lang := req.Header.Get("Accept-Language"); lang != "" {
			req.Header.Set("X-Atlas-Locale", lang)
		}
	}

	p.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		logger.Error("upstream reverse proxy error",
			"target", target.String(),
			"path", r.URL.Path,
			"error", err,
		)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`{"error":"upstream service unavailable"}`))
	}

	return &ReverseProxyHandler{
		target: target,
		proxy:  p,
	}, nil
}

func (h *ReverseProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Guard against unauthenticated callers trying to inject internal headers
	if !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
		r.Header.Del(middleware.HeaderUserID)
		r.Header.Del(middleware.HeaderUserEmail)
	}
	h.proxy.ServeHTTP(w, r)
}