package main

import (
	"context"
	"go-webapp/internal/config"
	"go-webapp/internal/handler"
	"go-webapp/internal/store"
	"go-webapp/internal/util"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	dataStore *store.Store
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Use configured frontend URL
			frontendUrl := config.FrontendURL

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

	// Initialize the data store
	dataStore = store.NewStore(24 * time.Hour) // 24 hour session duration

	googleOAuthConfig := config.GetGoogleOAuthConfig()
	githubOAuthConfig := config.GetGitHubOAuthConfig()

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

	port := config.Port
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
