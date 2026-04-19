package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("The .env file could not be charged")
	}
}

func main() {
	fmt.Println("Hello world")
}
