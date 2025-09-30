package handler

import (
	"encoding/json"
	"go-webapp/internal/store"
	"log"
	"net/http"
)

type UserResponse struct {
	Email           string `json:"email"`
	IsEmailVerified bool   `json:"isEmailVerified"`
	GoogleID        string `json:"googleId,omitempty"`
	GitHubID        string `json:"githubId,omitempty"`
	Name            string `json:"name,omitempty"`
	Picture         string `json:"picture,omitempty"`
	AuthProvider    string `json:"authProvider,omitempty"`
}

// Handle /me endpoint
func (h *Handler) HandleMeRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("🔍 Handle /me endpoint running")
	// Get user from context (set by middleware)
	user, ok := r.Context().Value("user").(store.User)
	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := UserResponse{
		Email:           user.Email,
		IsEmailVerified: user.IsEmailVerified,
		GoogleID:        user.GoogleID,
		GitHubID:        user.GitHubID,
		Name:            user.Name,
		Picture:         user.Picture,
		AuthProvider:    user.AuthProvider,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
