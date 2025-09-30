package handler

import (
	"context"
	"encoding/json"
	"go-webapp/internal/store"
	"go-webapp/internal/util"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

// Google OAuth handlers
func (h *Handler) HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	if h.googleOAuthConfig == nil {
		http.Error(w, "Google OAuth not configured", http.StatusInternalServerError)
		return
	}

	// Generate state parameter for CSRF protection
	state, err := util.GenerateSecureCode()
	if err != nil {
		log.Printf("Error generating state: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Store state with expiration (10 minutes)
	h.store.SetOAuthState(state, time.Now().Add(10*time.Minute), r.URL.Query().Get("return_to"))

	url := h.googleOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)

	log.Printf("Redirecting to Google OAuth: %s", url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	if h.googleOAuthConfig == nil {
		http.Error(w, "Google OAuth not configured", http.StatusInternalServerError)
		return
	}

	// Get authorization code and state from callback
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" {
		http.Error(w, "Authorization code not provided", http.StatusBadRequest)
		return
	}

	if state == "" {
		http.Error(w, "State parameter not provided", http.StatusBadRequest)
		return
	}

	// Validate state parameter for CSRF protection
	oauthState, exists := h.store.GetOAuthState(state)
	if !exists {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	if time.Now().After(oauthState.ExpiresAt) {
		h.store.DeleteOAuthState(state) // Clean up expired state
		http.Error(w, "State parameter has expired", http.StatusBadRequest)
		return
	}

	// Clean up used state
	h.store.DeleteOAuthState(state)

	// Exchange code for token
	token, err := h.googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Error exchanging code for token: %v", err)
		http.Error(w, "Failed to exchange authorization code", http.StatusInternalServerError)
		return
	}

	// Get user info from Google
	client := h.googleOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Printf("Error getting user info: %v", err)
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var googleUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		log.Printf("Error decoding user info: %v", err)
		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	// Create or update user
	user := store.User{
		Email:           googleUser.Email,
		IsEmailVerified: true, // Google emails are pre-verified
		GoogleID:        googleUser.ID,
		Name:            googleUser.Name,
		Picture:         googleUser.Picture,
		AuthProvider:    "google",
	}
	h.store.SetUser(googleUser.Email, user)

	// Create session
	sessionToken, err := h.store.CreateSession(googleUser.Email)
	if err != nil {
		log.Printf("Error creating session for user %s: %v", googleUser.Email, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	util.SetSessionCookie(w, sessionToken)

	log.Printf("User %s successfully logged in with Google", googleUser.Email)

	// Redirect directly to dashboard since session creation works correctly
	http.Redirect(w, r, util.GetFrontendReturnUrl(oauthState.ReturnTo), http.StatusSeeOther)
}
