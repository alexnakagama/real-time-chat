package config

import (
	"os"
)

var (
	JWTSecret = os.Getenv("JWT_SECRET")
)
