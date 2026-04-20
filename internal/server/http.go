package server

import (
	"encoding/json"
	"net/http"
	"real-time-chat/internal/auth"
)

type RegisterRequest struct {
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

func SetupRoutes() {
	http.Handle("/register", CORSMiddleware(RateLimitMiddleware(http.HandlerFunc(RegisterHandler))))
}
