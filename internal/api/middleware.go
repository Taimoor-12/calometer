package api

import (
	"calometer/internal/lib"
	"context"
	"net/http"

	"go.uber.org/zap"
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

func SetInitialTDEEMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the token from the context
		tokenStr, ok := r.Context().Value(TokenContextKey).(string)
		if !ok {
			log.Info(
				"token not found in context",
			)

			// Token is not present in context
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		userId, err := lib.ExtractUserIdFromToken(tokenStr)
		if err != nil {
			log.Info(
				"failed to get user id by token",
				zap.String("userId", userId.String()),
				zap.Error(err),
			)

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		exists, err := lib.DoesCalorieLogForUserExist(*userId)
		if err != nil {
			log.Info(
				"failed to determine user log's existence",
				zap.String("userId", userId.String()),
				zap.Error(err),
			)
		}

		if !*exists {
			bmr, err := lib.GetUserBmr(*userId)
			if err != nil {
				log.Info(
					"failed to get user bmr by id",
					zap.String("userId", userId.String()),
					zap.Error(err),
				)

				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if err := lib.SetInitialUserTdee(*userId, *bmr); err != nil {
				log.Info(
					"failed to set user's tdee by id",
					zap.String("userId", userId.String()),
					zap.Error(err),
				)

				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		ctx := context.WithValue(r.Context(), UserIdContextKey, *userId)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
