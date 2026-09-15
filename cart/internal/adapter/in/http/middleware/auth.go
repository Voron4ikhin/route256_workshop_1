package middleware

import (
	"net/http"
	"route256/cart/internal/authctx"
	"strings"
)

const bearerPrefix = "Bearer "

func Auth() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := strings.TrimPrefix(r.Header.Get("Authorization"), bearerPrefix)
			next.ServeHTTP(w, r.WithContext(authctx.WithToken(r.Context(), token)))
		})
	}
}
