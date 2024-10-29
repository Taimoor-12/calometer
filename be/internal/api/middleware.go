package api

import (
	"calometer/internal/lib"
	"context"
	"encoding/json"
	"net/http"
	"os"
)

type contextKey string

const TokenContextKey contextKey = "token"

type contextUserId string

const UserIdContextKey contextUserId = "userId"

func AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := Response{}
		resp.Code = make(map[int]string)

		// Extract the token from the cookie
		cookie, err := r.Cookie("token")
		if err != nil {
			// No cookie found
			resp.Code[http.StatusUnauthorized] = "Session expired. Please login again."
			json.NewEncoder(w).Encode(&resp)
			return
		}

		// Validate the JWT
		if err := lib.ValidateToken(cookie.Value); err != nil {
			// Token is invalid
			resp.Code[http.StatusUnauthorized] = "Invalid token. Please login again."
			json.NewEncoder(w).Encode(&resp)
			return
		}

		ctx := context.WithValue(r.Context(), TokenContextKey, cookie.Value)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("FE_URL"))
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		next.ServeHTTP(w, r)
	})
}
