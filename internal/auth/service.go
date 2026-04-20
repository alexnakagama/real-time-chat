package auth

import (
	"context"
	"errors"
	"real-time-chat/internal/db"
)

func RegisterUser(username, password string) error {
	hash, err := GeneratePasswordHash(password)
	if err != nil {
		return err
	}

	_, err = db.Pool.Exec(
		context.Background(),
		"INSERT INTO users (username, password_hash) VALUES ($1, $2)",
		username, hash,
	)
	return err
}

func AuthenticateUser(username, password string) (bool, error) {
	var passwordHash string
	err := db.Pool.QueryRow(
		context.Background(),
		"SELECT password_hash FROM users WHERE username = $1",
		username,
	).Scan(&passwordHash)

	if err != nil {
		return false, err
	}

	if !VerifyPasswordHash(password, passwordHash) {
		return false, errors.New("invalid credentials")
	}

	return true, nil
}
