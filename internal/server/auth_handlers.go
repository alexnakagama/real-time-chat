package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"real-time-chat/internal/auth"
	"real-time-chat/internal/db"
	"time"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
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

// ForgotPasswordHandler godoc
// @Summary      Forgot Password
// @Description  Sends a password reset link to the user's email
// @Tags         auth
// @Accept       json
// @Produce      plain
// @Param        request body ForgotPasswordRequest true "User email"
// @Success      200 {string} string "reset link sent"
// @Failure      400 {string} string "invalid request"
// @Failure      500 {string} string "error sending reset link"
// @Router       /forgot-password [post]
func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user, err := db.GetUserByEmail(req.Email)
	if err != nil {
		w.Write([]byte("reset link sent"))
		return
	}

	token, err := auth.GenerateResetToken()
	if err != nil {
		http.Error(w, "error generating token", http.StatusInternalServerError)
		return
	}

	err = db.SavePasswordResetToken(user.ID, token, time.Now().Add(1*time.Hour))
	if err != nil {
		http.Error(w, "error saving token", http.StatusInternalServerError)
		return
	}

	resetLink := fmt.Sprintf("https://tuapp.com/reset-password?token=%s", token)
	fmt.Println("Reset link:", resetLink)
	// In prod here you send the real one

	w.Write([]byte("reset link sent"))
}

// ResetPasswordHandler godoc
// @Summary      Reset Password
// @Description  Resets the user password using a valid reset token
// @Tags         auth
// @Accept       json
// @Produce      plain
// @Param        request body ResetPasswordRequest true "Reset password data"
// @Success      200 {string} string "password reset successful"
// @Failure      400 {string} string "invalid request"
// @Failure      401 {string} string "invalid or expired token"
// @Failure      500 {string} string "could not reset password"
// @Router       /reset-password [post]
func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := auth.ResetPassword(req.Token, req.NewPassword)
	if err != nil {
		http.Error(w, "invalid or expired token", http.StatusUnauthorized)
		return
	}

	w.Write([]byte("password reset successful"))
}
