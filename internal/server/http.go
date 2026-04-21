package server

import (
	"net/http"
	"real-time-chat/internal/chat"
)

var Hub = chat.NewHub()

func SetupRoutes() {
	http.Handle("/register", CORSMiddleware(RateLimitMiddleware(http.HandlerFunc(RegisterHandler))))
	http.Handle("/login", CORSMiddleware(RateLimitMiddleware(http.HandlerFunc(LoginHandler))))
	http.Handle("/forgot-password", CORSMiddleware(RateLimitMiddleware(http.HandlerFunc(ForgotPasswordHandler))))
	http.Handle("/reset-password", CORSMiddleware(RateLimitMiddleware(http.HandlerFunc(ResetPasswordHandler))))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(Hub, w, r)
	})
}
