package middleware

import (
	"context"
	"net/http"

	httphandler "github.com/aras-services/aras-auth/internal/delivery/http"
	"github.com/aras-services/aras-auth/internal/domain"
)

type AuthMiddleware struct {
	tokenService domain.TokenService
}

func NewAuthMiddleware(tokenService domain.TokenService) *AuthMiddleware {
	return &AuthMiddleware{
		tokenService: tokenService,
	}
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			httphandler.WriteUnauthorized(w, "Authorization header required")
			return
		}

		// Extract token from "Bearer <token>" format
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			httphandler.WriteUnauthorized(w, "Invalid authorization header format")
			return
		}
		token := authHeader[7:]

		// Validate token
		claims, err := m.tokenService.ValidateAccessToken(token)
		if err != nil {
			httphandler.WriteUnauthorized(w, "Invalid or expired token")
			return
		}

		// Add user information to context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID.String())
		ctx = context.WithValue(ctx, "user_email", claims.Email)
		ctx = context.WithValue(ctx, "token_claims", claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Extract token from "Bearer <token>" format
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			next.ServeHTTP(w, r)
			return
		}
		token := authHeader[7:]

		// Validate token
		claims, err := m.tokenService.ValidateAccessToken(token)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		// Add user information to context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID.String())
		ctx = context.WithValue(ctx, "user_email", claims.Email)
		ctx = context.WithValue(ctx, "token_claims", claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
