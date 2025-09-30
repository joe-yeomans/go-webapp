package handler

import (
	"encoding/json"
	"fmt"
	"go-webapp/internal/store"
	"go-webapp/internal/util"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type CodeRequest struct {
	Email string `json:"email"`
}

type CodeResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// validateReturnToURL validates that the return_to URL is safe to redirect to
// It ensures the URL is a relative path within the application
func validateReturnToURL(returnTo string) bool {
	if returnTo == "" {
		return true // Empty is valid, will use default
	}

	// Must start with / (relative path)
	if !strings.HasPrefix(returnTo, "/") {
		return false
	}

	// Must not contain protocol or host (no http://, https://, //, etc.)
	if strings.Contains(returnTo, "://") || strings.HasPrefix(returnTo, "//") {
		return false
	}

	// Must not contain query parameters with external URLs
	if strings.Contains(returnTo, "?") && strings.Contains(returnTo, "http") {
		return false
	}

	return true
}

func (h *Handler) HandleCodeRequest(w http.ResponseWriter, r *http.Request) {
	var req CodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Generate a secure 6-digit code
	code, err := util.GenerateSecureCode()
	if err != nil {
		log.Printf("Error generating code: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Println("Code generated:", code)

	loginCode := store.LoginCode{
		Email:     req.Email,
		Code:      code,
		ExpiresAt: time.Now().Add(10 * time.Minute), // 10 minutes expiry
	}

	h.store.SetLoginCode(req.Email, loginCode)

	response := CodeResponse{
		Success: true,
		Message: "Code sent successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) HandleVerifyCodeRequest(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	email := r.FormValue("email")
	returnTo := r.FormValue("return_to")

	// Frontend URL - use environment variable or default to localhost:3000
	frontendUrl := util.GetEnvString("FRONTEND_URL", "http://localhost:3000")

	// Validate required fields
	if code == "" {
		http.Redirect(w, r, fmt.Sprintf("%s/verify?email=%s&error=%s", frontendUrl, url.QueryEscape(email), url.QueryEscape("Code is required")), http.StatusSeeOther)
		return
	}
	if email == "" {
		http.Redirect(w, r, fmt.Sprintf("%s/verify?error=%s", frontendUrl, url.QueryEscape("Email is required")), http.StatusSeeOther)
		return
	}

	// Get the stored login code
	loginCode, exists := h.store.GetLoginCode(email)
	if !exists {
		http.Redirect(w, r, fmt.Sprintf("%s/verify?email=%s&error=%s", frontendUrl, url.QueryEscape(email), url.QueryEscape("No verification code found for this email")), http.StatusSeeOther)
		return
	}

	// Check if code has expired
	if time.Now().After(loginCode.ExpiresAt) {
		h.store.DeleteLoginCode(email) // Clean up expired code
		http.Redirect(w, r, fmt.Sprintf("%s/verify?email=%s&error=%s", frontendUrl, url.QueryEscape(email), url.QueryEscape("Verification code has expired")), http.StatusSeeOther)
		return
	}

	// Verify the code
	if code != loginCode.Code {
		http.Redirect(w, r, fmt.Sprintf("%s/verify?email=%s&error=%s", frontendUrl, url.QueryEscape(email), url.QueryEscape("Invalid verification code")), http.StatusSeeOther)
		return
	}

	// Code is valid - clean up the used code
	h.store.DeleteLoginCode(email)

	// Create or update user as verified
	user := store.User{
		Email:           email,
		IsEmailVerified: true,
	}
	h.store.SetUser(email, user)

	// Create a new session
	sessionToken, err := h.store.CreateSession(email)
	if err != nil {
		log.Printf("Error creating session for user %s: %v", email, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("User %s successfully verified with code %s", email, code)

	util.SetSessionCookie(w, sessionToken)

	// Validate return_to URL to prevent open redirects
	if !validateReturnToURL(returnTo) {
		log.Printf("Invalid return_to URL provided: %s", returnTo)
		returnTo = "/start" // Use default if invalid
	}

	// Determine redirect destination
	redirectPath := "/start" // Default
	if returnTo != "" && validateReturnToURL(returnTo) {
		redirectPath = returnTo
	}

	http.Redirect(w, r, fmt.Sprintf("%s%s", frontendUrl, redirectPath), http.StatusSeeOther)
}
