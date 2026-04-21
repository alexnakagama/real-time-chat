package auth

import (
	"context"
	"fmt"
	"real-time-chat/internal/db"
)

func RegisterUser(username, password string, email string) error {
	hash, err := GeneratePasswordHash(password)
	if err != nil {
		return err
	}

	_, err = db.Pool.Exec(
		context.Background(),
		"INSERT INTO users (username, password_hash, email) VALUES ($1, $2, $3)",
		username, hash, email,
	)
	return err
}

func ResetPassword(token, newPassword string) error {
	userID, err := db.GetUserIDByResetToken(token)
	if err != nil {
		return fmt.Errorf("invalid or expired token")
	}

	hash, err := GeneratePasswordHash(newPassword)
	if err != nil {
		return fmt.Errorf("could not hash password: %w", err)
	}

	err = db.UpdateUserPassword(userID, hash)
	if err != nil {
		return fmt.Errorf("could not update password: %w", err)
	}

	_ = db.DeletePasswordResetToken(token)

	return nil
}
