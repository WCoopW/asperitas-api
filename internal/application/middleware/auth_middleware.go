package middleware

import (
	"context"
	"net/http"
	"strings"

	"reddit/internal/domain/auth"
	"reddit/pkg/helpers"
)

const ctxUserID = "user_id"

type AuthMiddleware struct {
	validator auth.TokenValidator
}

func NewAuthMiddleware(validator auth.TokenValidator) *AuthMiddleware {
	return &AuthMiddleware{validator: validator}
}

func (m *AuthMiddleware) WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := strings.TrimSpace(r.Header.Get("Authorization"))
		if raw == "" {
			helpers.WriteError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		token := raw
		if len(raw) > 7 && strings.EqualFold(raw[:7], "Bearer ") {
			token = strings.TrimSpace(raw[7:])
		}
		userID, err := m.validator.ValidateToken(token)
		if err != nil {
			helpers.WriteError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		ctx := context.WithValue(r.Context(), ctxUserID, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
