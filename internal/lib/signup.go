package lib

import (
	"calometer/internal/db"
	"context"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func DoesUserExists(email string) (*bool, error) {
	qStr := `
		SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE email = $1
		)`

	var doesExist bool
	if err := db.GetPool().QueryRow(context.Background(), qStr, email).Scan(&doesExist); err != nil {
		return nil, err
	}

	return &doesExist, nil
}
