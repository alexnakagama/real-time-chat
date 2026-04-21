package db

import (
	"context"
	"fmt"

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
