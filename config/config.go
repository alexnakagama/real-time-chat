package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

var (
	Database  DBConfig
	JWTSecret string
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Enviroment variables could not be charged")
	}
	JWTSecret = os.Getenv("JWT_SECRET")

	Database = DBConfig{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Name:     os.Getenv("DB_NAME"),
	}
}
