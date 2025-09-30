package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"go-webapp/internal/store"
	"go-webapp/internal/util"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

// GitHub OAuth handlers
func (h *Handler) HandleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	if h.githubOAuthConfig == nil {
		http.Error(w, "GitHub OAuth not configured", http.StatusInternalServerError)
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

	url := h.githubOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)

	log.Printf("Redirecting to GitHub OAuth: %s", url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) HandleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	if h.githubOAuthConfig == nil {
		http.Error(w, "GitHub OAuth not configured", http.StatusInternalServerError)
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
	token, err := h.githubOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Error exchanging code for token: %v", err)
		http.Error(w, "Failed to exchange authorization code", http.StatusInternalServerError)
		return
	}

	// Get user info from GitHub
	client := h.githubOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		log.Printf("Error getting user info: %v", err)
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var githubUser struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		log.Printf("Error decoding user info: %v", err)
		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	// If email is not public, get it from GitHub's email endpoint
	if githubUser.Email == "" {
		emailResp, err := client.Get("https://api.github.com/user/emails")
		if err == nil {
			defer emailResp.Body.Close()

			var emails []struct {
				Email   string `json:"email"`
				Primary bool   `json:"primary"`
			}

			if err := json.NewDecoder(emailResp.Body).Decode(&emails); err == nil {
				for _, email := range emails {
					if email.Primary {
						githubUser.Email = email.Email
						break
					}
				}
				// If no primary email found, use the first one
				if githubUser.Email == "" && len(emails) > 0 {
					githubUser.Email = emails[0].Email
				}
			}
		}
	}

	if githubUser.Email == "" {
		http.Error(w, "Unable to get email from GitHub", http.StatusInternalServerError)
		return
	}

	// Use login as name if name is not available
	if githubUser.Name == "" {
		githubUser.Name = githubUser.Login
	}

	// Create or update user
	user := store.User{
		Email:           githubUser.Email,
		IsEmailVerified: true, // GitHub emails are pre-verified
		GitHubID:        fmt.Sprintf("%d", githubUser.ID),
		Name:            githubUser.Name,
		Picture:         githubUser.AvatarURL,
		AuthProvider:    "github",
	}
	h.store.SetUser(githubUser.Email, user)

	// Create session
	sessionToken, err := h.store.CreateSession(githubUser.Email)
	if err != nil {
		log.Printf("Error creating session for user %s: %v", githubUser.Email, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	util.SetSessionCookie(w, sessionToken)

	log.Printf("User %s successfully logged in with GitHub", githubUser.Email)
	http.Redirect(w, r, util.GetFrontendReturnUrl(oauthState.ReturnTo), http.StatusSeeOther)
}
