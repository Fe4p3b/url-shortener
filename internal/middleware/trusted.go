package middleware

import (
	"net/http"
)

type TrustedNetworksOnlyMiddleware struct {
	IPs map[string]struct{}
}

func NewTrustedNetworksOnlyMiddleware(IPs []string) *TrustedNetworksOnlyMiddleware {
	m := make(map[string]struct{})
	for _, v := range IPs {
		m[v] = struct{}{}
	}
	return &TrustedNetworksOnlyMiddleware{IPs: m}
}

func (t *TrustedNetworksOnlyMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xRealIP := r.Header.Get("X-Real-IP")
		if _, ok := t.IPs[xRealIP]; !ok {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
