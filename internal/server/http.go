package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"real-time-chat/internal/auth"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterHandler godoc
// @Summary      Register User
// @Description  Registers a new user
// @Tags         auth
// @Accept       json
// @Produce      plain
// @Param        request body RegisterRequest true "Data of register"
// @Success      201 {string} string "user registered"
// @Failure      400 {string} string "invalid request"
// @Failure      500 {string} string "error registering user"
// @Router       /register [post]
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := auth.RegisterUser(req.Username, req.Password)
	if err != nil {
		http.Error(w, "error registering user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("user registered"))
}

// LoginHandler godoc
// @Summary      Login User
// @Description  Authenticates a user and returns a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login credentials"
// @Success      200 {object} map[string]string "token"
// @Failure      400 {string} string "invalid credentials"
// @Failure      401 {string} string "invalid credentials"
// @Failure      500 {string} string "could not generate token"
// @Router       /login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	user, err := auth.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateJWT(fmt.Sprintf("%d", user.ID), user.Username)
	if err != nil {
		http.Error(w, "could not generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func SetupRoutes() {
	http.Handle("/register", CORSMiddleware(RateLimitMiddleware(http.HandlerFunc(RegisterHandler))))
	http.Handle("/login", CORSMiddleware(RateLimitMiddleware(http.HandlerFunc(LoginHandler))))
}
