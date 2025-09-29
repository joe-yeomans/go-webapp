package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"go-webapp/internal/store"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	dataStore         *store.Store
	googleOAuthConfig *oauth2.Config
	githubOAuthConfig *oauth2.Config
	oauthStates       = make(map[string]time.Time) // Store OAuth states with expiration
)

// generateSecureCode generates a cryptographically secure 6-digit code
func generateSecureCode() (string, error) {
	// Generate a random number between 0 and 999999
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	// Format as 6-digit string with leading zeros
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func setSessionCookie(w http.ResponseWriter, sessionToken string) {
	// Use secure cookies only in production (HTTPS)
	frontendUrl := getEnvString("FRONTEND_URL", "http://localhost:3000")
	cookieDomain := getEnvString("COOKIE_DOMAIN", "localhost")
	isSecure := strings.HasPrefix(frontendUrl, "https://")

	// Get session time from data store
	sessionTime := dataStore.GetSessionTime()

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionToken,
		Path:     "/",
		Domain:   cookieDomain,
		MaxAge:   int(sessionTime.Seconds()),
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteLaxMode,
	})
}

func clearSessionCookie(w http.ResponseWriter) {
	// Use secure cookies only in production (HTTPS)
	frontendUrl := getEnvString("FRONTEND_URL", "http://localhost:3000")
	cookieDomain := getEnvString("COOKIE_DOMAIN", "localhost")
	isSecure := strings.HasPrefix(frontendUrl, "https://")

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		Domain:   cookieDomain,
		MaxAge:   -1, // Delete the cookie
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteLaxMode,
	})
}

func loadConfig() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Initialize session time from environment variable
	sessionTimeMinutes := getEnvInt("SESSION_TIME_MINUTES", 10)
	sessionTime := time.Duration(sessionTimeMinutes) * time.Minute

	// Initialize data store with session time
	dataStore = store.NewStore(sessionTime)

	// Get server configuration
	port := getEnvString("PORT", "8080")
	frontendUrl := getEnvString("FRONTEND_URL", "http://localhost:3000")
	cookieDomain := getEnvString("COOKIE_DOMAIN", "localhost")

	// Google OAuth configuration
	googleClientId := getEnvString("GOOGLE_CLIENT_ID", "")
	googleClientSecret := getEnvString("GOOGLE_CLIENT_SECRET", "")
	googleRedirectUrl := getEnvString("GOOGLE_REDIRECT_URL", fmt.Sprintf("http://localhost:%s/api/auth/google/callback", port))

	if googleClientId != "" && googleClientSecret != "" {
		googleOAuthConfig = &oauth2.Config{
			ClientID:     googleClientId,
			ClientSecret: googleClientSecret,
			RedirectURL:  googleRedirectUrl,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		}
		log.Println("Google OAuth configured")
	} else {
		log.Println("Warning: Google OAuth not configured - missing GOOGLE_CLIENT_ID or GOOGLE_CLIENT_SECRET")
	}

	// GitHub OAuth configuration
	githubClientId := getEnvString("GITHUB_CLIENT_ID", "")
	githubClientSecret := getEnvString("GITHUB_CLIENT_SECRET", "")
	githubRedirectUrl := getEnvString("GITHUB_REDIRECT_URL", fmt.Sprintf("http://localhost:%s/api/auth/github/callback", port))

	if githubClientId != "" && githubClientSecret != "" {
		githubOAuthConfig = &oauth2.Config{
			ClientID:     githubClientId,
			ClientSecret: githubClientSecret,
			RedirectURL:  githubRedirectUrl,
			Scopes: []string{
				"user:email", // Access user email
			},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://github.com/login/oauth/authorize",
				TokenURL: "https://github.com/login/oauth/access_token",
			},
		}
		log.Println("GitHub OAuth configured")
	} else {
		log.Println("Warning: GitHub OAuth not configured - missing GITHUB_CLIENT_ID or GITHUB_CLIENT_SECRET")
	}

	log.Printf("Configuration loaded - Frontend URL: %s, Cookie Domain: %s", frontendUrl, cookieDomain)
}

// Helper functions for environment variables
func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func main() {
	// Load configuration from .env file
	loadConfig()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// Routes
	r.Route("/api", func(r chi.Router) {
		r.Post("/code", handleCodeRequest)
		r.Post("/verify", handleVerifyCodeRequest)
		r.Get("/me", sessionMiddleware(handleMeRequest))
		r.Post("/logout", sessionMiddleware(handleLogoutRequest))

		// Google OAuth routes
		r.Get("/auth/google", handleGoogleLogin)
		r.Get("/auth/google/callback", handleGoogleCallback)

		// GitHub OAuth routes
		r.Get("/auth/github", handleGitHubLogin)
		r.Get("/auth/github/callback", handleGitHubCallback)
	})

	port := getEnvString("PORT", "8080")
	log.Printf("Server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

type CodeRequest struct {
	Email string `json:"email"`
}

type CodeResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func handleCodeRequest(w http.ResponseWriter, r *http.Request) {
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
	code, err := generateSecureCode()
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

	dataStore.SetLoginCode(req.Email, loginCode)

	response := CodeResponse{
		Success: true,
		Message: "Code sent successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleVerifyCodeRequest(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	email := r.FormValue("email")

	// Frontend URL - use environment variable or default to localhost:3000
	frontendUrl := getEnvString("FRONTEND_URL", "http://localhost:3000")

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
	loginCode, exists := dataStore.GetLoginCode(email)
	if !exists {
		http.Redirect(w, r, fmt.Sprintf("%s/verify?email=%s&error=%s", frontendUrl, url.QueryEscape(email), url.QueryEscape("No verification code found for this email")), http.StatusSeeOther)
		return
	}

	// Check if code has expired
	if time.Now().After(loginCode.ExpiresAt) {
		dataStore.DeleteLoginCode(email) // Clean up expired code
		http.Redirect(w, r, fmt.Sprintf("%s/verify?email=%s&error=%s", frontendUrl, url.QueryEscape(email), url.QueryEscape("Verification code has expired")), http.StatusSeeOther)
		return
	}

	// Verify the code
	if code != loginCode.Code {
		http.Redirect(w, r, fmt.Sprintf("%s/verify?email=%s&error=%s", frontendUrl, url.QueryEscape(email), url.QueryEscape("Invalid verification code")), http.StatusSeeOther)
		return
	}

	// Code is valid - clean up the used code
	dataStore.DeleteLoginCode(email)

	// Create or update user as verified
	user := store.User{
		Email:           email,
		IsEmailVerified: true,
	}
	dataStore.SetUser(email, user)

	// Create a new session
	sessionToken, err := dataStore.CreateSession(email)
	if err != nil {
		log.Printf("Error creating session for user %s: %v", email, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("User %s successfully verified with code %s", email, code)

	setSessionCookie(w, sessionToken)

	http.Redirect(w, r, fmt.Sprintf("%s/start", frontendUrl), http.StatusSeeOther)
}

// Session middleware to validate authentication
func sessionMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("🔍 Session middleware running")

		// Get session cookie
		cookie, err := r.Cookie("session")
		if err != nil {
			log.Println("❌ No session cookie")
			http.Error(w, "Unauthorized: No session cookie", http.StatusUnauthorized)
			return
		}

		// Get session from store
		session, exists := dataStore.GetSession(cookie.Value)
		if !exists {
			log.Println("❌ Invalid session token")
			clearSessionCookie(w)
			http.Error(w, "Unauthorized: Invalid session", http.StatusUnauthorized)
			return
		}

		// Check if session has expired
		if !dataStore.IsSessionValid(cookie.Value) {
			log.Println("❌ Session expired")
			// Clean up expired session
			dataStore.DeleteSession(cookie.Value)
			// Clear the expired cookie
			clearSessionCookie(w)
			http.Error(w, "Unauthorized: Session expired", http.StatusUnauthorized)
			return
		}

		// Get user from store
		user, exists := dataStore.GetUser(session.UserEmail)
		if !exists {
			log.Println("❌ User not found")
			http.Error(w, "Unauthorized: User not found", http.StatusUnauthorized)
			return
		}

		// Check if user is verified
		if !user.IsEmailVerified {
			log.Println("❌ Email not verified")
			http.Error(w, "Unauthorized: Email not verified", http.StatusUnauthorized)
			return
		}

		// Add user to request context for use in handlers
		ctx := context.WithValue(r.Context(), "user", user)
		r = r.WithContext(ctx)

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// User response type
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
func handleMeRequest(w http.ResponseWriter, r *http.Request) {
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

// Google OAuth handlers
func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	if googleOAuthConfig == nil {
		http.Error(w, "Google OAuth not configured", http.StatusInternalServerError)
		return
	}

	// Generate state parameter for CSRF protection
	state, err := generateSecureCode()
	if err != nil {
		log.Printf("Error generating state: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Store state with expiration (10 minutes)
	oauthStates[state] = time.Now().Add(10 * time.Minute)

	url := googleOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)

	log.Printf("Redirecting to Google OAuth: %s", url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	if googleOAuthConfig == nil {
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
	expiration, exists := oauthStates[state]
	if !exists {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	if time.Now().After(expiration) {
		delete(oauthStates, state) // Clean up expired state
		http.Error(w, "State parameter has expired", http.StatusBadRequest)
		return
	}

	// Clean up used state
	delete(oauthStates, state)

	// Exchange code for token
	token, err := googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Error exchanging code for token: %v", err)
		http.Error(w, "Failed to exchange authorization code", http.StatusInternalServerError)
		return
	}

	// Get user info from Google
	client := googleOAuthConfig.Client(context.Background(), token)
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
	dataStore.SetUser(googleUser.Email, user)

	// Create session
	sessionToken, err := dataStore.CreateSession(googleUser.Email)
	if err != nil {
		log.Printf("Error creating session for user %s: %v", googleUser.Email, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	setSessionCookie(w, sessionToken)

	// Redirect to dashboard
	frontendUrl := getEnvString("FRONTEND_URL", "http://localhost:3000")

	log.Printf("User %s successfully logged in with Google", googleUser.Email)

	// Redirect directly to dashboard since session creation works correctly
	http.Redirect(w, r, fmt.Sprintf("%s/dashboard", frontendUrl), http.StatusSeeOther)
}

// GitHub OAuth handlers
func handleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	if githubOAuthConfig == nil {
		http.Error(w, "GitHub OAuth not configured", http.StatusInternalServerError)
		return
	}

	// Generate state parameter for CSRF protection
	state, err := generateSecureCode()
	if err != nil {
		log.Printf("Error generating state: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Store state with expiration (10 minutes)
	oauthStates[state] = time.Now().Add(10 * time.Minute)

	url := githubOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)

	log.Printf("Redirecting to GitHub OAuth: %s", url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	if githubOAuthConfig == nil {
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
	expiration, exists := oauthStates[state]
	if !exists {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	if time.Now().After(expiration) {
		delete(oauthStates, state) // Clean up expired state
		http.Error(w, "State parameter has expired", http.StatusBadRequest)
		return
	}

	// Clean up used state
	delete(oauthStates, state)

	// Exchange code for token
	token, err := githubOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Error exchanging code for token: %v", err)
		http.Error(w, "Failed to exchange authorization code", http.StatusInternalServerError)
		return
	}

	// Get user info from GitHub
	client := githubOAuthConfig.Client(context.Background(), token)
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
	dataStore.SetUser(githubUser.Email, user)

	// Create session
	sessionToken, err := dataStore.CreateSession(githubUser.Email)
	if err != nil {
		log.Printf("Error creating session for user %s: %v", githubUser.Email, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	setSessionCookie(w, sessionToken)

	// Redirect to dashboard
	frontendUrl := getEnvString("FRONTEND_URL", "http://localhost:3000")

	log.Printf("User %s successfully logged in with GitHub", githubUser.Email)
	http.Redirect(w, r, fmt.Sprintf("%s/dashboard", frontendUrl), http.StatusSeeOther)
}

// Logout handler
func handleLogoutRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("🔍 Handle logout endpoint running")

	// Get user from context (set by middleware)
	user, ok := r.Context().Value("user").(store.User)
	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get session cookie
	cookie, err := r.Cookie("session")
	if err == nil {
		// Delete session from store
		dataStore.DeleteSession(cookie.Value)
		log.Printf("Session deleted for user %s", user.Email)
	}

	clearSessionCookie(w)

	// Get frontend URL
	frontendUrl := getEnvString("FRONTEND_URL", "http://localhost:3000")

	log.Printf("User %s successfully logged out", user.Email)

	// Redirect to frontend home page
	http.Redirect(w, r, frontendUrl, http.StatusSeeOther)
}
