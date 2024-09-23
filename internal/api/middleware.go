package api

import (
	"calometer/internal/lib"
	"context"
	"net/http"
)

type contextKey string

const TokenContextKey contextKey = "token"

type contextUserId string

const UserIdContextKey contextUserId = "userId"

func AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for JWT in the request
		tokenStr := lib.ExtractTokenFromHeader(r)
		if tokenStr == "" {
			// No token found in the request header
			http.Error(w, "User not authenticated. Please provide a valid token.", http.StatusUnauthorized)
			return
		}

		// Validate the JWT
		if err := lib.ValidateToken(tokenStr); err != nil {
			// Token is invalid
			http.Error(w, "Invalid token. Please authenticate again.", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), TokenContextKey, tokenStr)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
