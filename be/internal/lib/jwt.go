package lib

import (
	"calometer/internal/db"
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func ExtractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// Expecting the header in the format "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1]
	}

	return ""
}

func ValidateToken(tokenStr string) error {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Validate the algorithm used to sign the token
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKey
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return err
	}

	return nil
}

func GenerateJWT(userId uuid.UUID, username string) (string, error) {
	claims := jwt.MapClaims{
		"u_id":     userId,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func GetHashedPass(username string) (string, error) {
	qStr := `
		SELECT password_hash
		FROM users
		WHERE username = $1`

	var passwordHash string
	if err := db.GetPool().QueryRow(
		context.Background(),
		qStr,
		username,
	).Scan(&passwordHash); err != nil {
		return "", err
	}

	return passwordHash, nil
}

func CheckPasswordValidity(password, passwordHash string) error {
	if err := bcrypt.CompareHashAndPassword(
		[]byte(passwordHash),
		[]byte(password),
	); err != nil {
		return err
	}

	return nil
}

func ExtractUserIdFromToken(tokenStr string) (*uuid.UUID, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Check if the signing method is the expected one
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKey
		}
		// Return the secret used to sign the token
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	// Extract the claims and retrieve the userId
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userIdStr, ok := claims["u_id"].(string); ok {
			userId, err := uuid.Parse(userIdStr)
			if err != nil {
				return nil, err
			}

			return &userId, nil
		}
	}

	return nil, errors.New("userId not found in token")
}
