package db

import (
	"context"
	"fmt"

	"real-time-chat/config"

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
