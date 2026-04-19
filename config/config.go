package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var JWTSecret string

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Enviroment variables could not be charged")
	}
	JWTSecret = os.Getenv("JWT_SECRET")
}
