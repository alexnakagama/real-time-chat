package db

import (
	"context"
	"fmt"
	"time"

	"real-time-chat/config"
	"real-time-chat/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect() error {
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name,
	)
	var err error
	Pool, err = pgxpool.New(context.Background(), url)
	return err
}

func Close() {
	if Pool != nil {
		Pool.Close()
	}
}

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	query := "SELECT id, username, password_hash FROM users WHERE username = $1"
	err := Pool.QueryRow(context.Background(), query, username).Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := "SELECT id, username, password_hash FROM users WHERE email = $1"
	err := Pool.QueryRow(context.Background(), query, email).Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func SavePasswordResetToken(userID int, token string, expiresAt time.Time) error {
	query := `
        INSERT INTO password_resets (user_id, token, expires_at)
        VALUES ($1, $2, $3)
        ON CONFLICT (user_id) DO UPDATE
        SET token = EXCLUDED.token, expires_at = EXCLUDED.expires_at
    `
	_, err := Pool.Exec(context.Background(), query, userID, token, expiresAt)
	return err
}

func GetUserIDByResetToken(token string) (int, error) {
	var userID int
	var expiresAt time.Time
	query := "SELECT user_id, expires_at FROM password_resets WHERE token = $1"
	err := Pool.QueryRow(context.Background(), query, token).Scan(&userID, &expiresAt)
	if err != nil {
		return 0, err
	}
	if time.Now().After(expiresAt) {
		return 0, fmt.Errorf("token expired")
	}
	return userID, nil
}

func DeletePasswordResetToken(token string) error {
	query := "DELETE FROM password_resets WHERE token = $1"
	_, err := Pool.Exec(context.Background(), query, token)
	return err
}

func UpdateUserPassword(userID int, newPasswordHash string) error {
	query := "UPDATE users SET password_hash = $1 WHERE id = $2"
	_, err := Pool.Exec(context.Background(), query, newPasswordHash, userID)
	return err
}
