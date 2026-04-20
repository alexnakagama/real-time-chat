// @title           Real Time Chat API
// @version         1.0
// @description     API for real-time chat application
// @host            localhost:8000
// @BasePath        /
package main

import (
	"log"
	"net/http"
	"real-time-chat/config"
	"real-time-chat/internal/db"
	"real-time-chat/internal/server"

	_ "real-time-chat/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	config.LoadEnv()

	if err := db.Connect(); err != nil {
		log.Fatal("Error conecting to database:", err)
	}
	defer db.Close()

	server.SetupRoutes()

	http.Handle("/swagger/", httpSwagger.WrapHandler)

	log.Println("Server listening on port :8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal("Error initializing server:", err)
	}
}
