package main

import (
	"fmt"
	"real-time-chat/config"
)

func main() {
	config.LoadEnv()
	fmt.Println("Hello world")
}
