package main

import (
	"context"
	"fmt"
	"go-webapp/internal/handler"
	"go-webapp/internal/store"
	"go-webapp/internal/util"
	"log"
	"net/http"
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

func loadConfig() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Initialize session time from environment variable
	sessionTimeMinutes := util.GetEnvInt("SESSION_TIME_MINUTES", 10)
	sessionTime := time.Duration(sessionTimeMinutes) * time.Minute

	// Initialize data store with session time
	dataStore = store.NewStore(sessionTime)

	// Get server configuration
	port := util.GetEnvString("PORT", "8080")
	frontendUrl := util.GetEnvString("FRONTEND_URL", "http://localhost:3000")
	cookieDomain := util.GetEnvString("COOKIE_DOMAIN", "localhost")

	// Google OAuth configuration
	googleClientId := util.GetEnvString("GOOGLE_CLIENT_ID", "")
	googleClientSecret := util.GetEnvString("GOOGLE_CLIENT_SECRET", "")
	googleRedirectUrl := util.GetEnvString("GOOGLE_REDIRECT_URL", fmt.Sprintf("http://localhost:%s/api/auth/google/callback", port))

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
	githubClientId := util.GetEnvString("GITHUB_CLIENT_ID", "")
	githubClientSecret := util.GetEnvString("GITHUB_CLIENT_SECRET", "")
	githubRedirectUrl := util.GetEnvString("GITHUB_REDIRECT_URL", fmt.Sprintf("http://localhost:%s/api/auth/github/callback", port))

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
			// Get the frontend URL from environment or default to localhost:3000
			frontendUrl := util.GetEnvString("FRONTEND_URL", "http://localhost:3000")

			w.Header().Set("Access-Control-Allow-Origin", frontendUrl)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token, Cookie")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	handler := handler.NewHandler(dataStore, googleOAuthConfig, githubOAuthConfig)

	// Routes
	r.Route("/api", func(r chi.Router) {
		r.Post("/code", handler.HandleCodeRequest)
		r.Post("/verify", handler.HandleVerifyCodeRequest)
		r.Get("/me", sessionMiddleware(handler.HandleMeRequest))
		r.Post("/logout", sessionMiddleware(handler.HandleLogoutRequest))

		// Google OAuth routes
		r.Get("/auth/google", handler.HandleGoogleLogin)
		r.Get("/auth/google/callback", handler.HandleGoogleCallback)

		// GitHub OAuth routes
		r.Get("/auth/github", handler.HandleGitHubLogin)
		r.Get("/auth/github/callback", handler.HandleGitHubCallback)

		// Product API routes (protected)
		r.Route("/products", func(r chi.Router) {
			r.Get("/", sessionMiddleware(handler.HandleGetProducts))
			r.Get("/{id}", sessionMiddleware(handler.HandleGetProduct))
			r.Get("/category/{category}", sessionMiddleware(handler.HandleGetProductsByCategory))
		})
	})

	port := util.GetEnvString("PORT", "8080")
	log.Printf("Server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
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
			util.ClearSessionCookie(w)
			http.Error(w, "Unauthorized: Invalid session", http.StatusUnauthorized)
			return
		}

		// Check if session has expired
		if !dataStore.IsSessionValid(cookie.Value) {
			log.Println("❌ Session expired")
			// Clean up expired session
			dataStore.DeleteSession(cookie.Value)
			// Clear the expired cookie
			util.ClearSessionCookie(w)
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
