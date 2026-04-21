package auth

import (
	"context"
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
