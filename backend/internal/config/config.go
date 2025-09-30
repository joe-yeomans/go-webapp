package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	Port               string
	FrontendURL        string
	CookieDomain       string
	SessionTimeMinutes int

	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string

	GitHubClientID     string
	GitHubClientSecret string
	GitHubRedirectURL  string
)

func init() {
	_ = godotenv.Load()

	Port = getEnvString("PORT", "8080")
	FrontendURL = getEnvString("FRONTEND_URL", "http://localhost:3000")
	CookieDomain = getEnvString("COOKIE_DOMAIN", "localhost")
	SessionTimeMinutes = getEnvInt("SESSION_TIME_MINUTES", 10)

	GoogleClientID = getEnvString("GOOGLE_CLIENT_ID", "")
	GoogleClientSecret = getEnvString("GOOGLE_CLIENT_SECRET", "")
	GoogleRedirectURL = getEnvString("GOOGLE_REDIRECT_URL", fmt.Sprintf("http://localhost:%s/api/auth/google/callback", Port))

	GitHubClientID = getEnvString("GITHUB_CLIENT_ID", "")
	GitHubClientSecret = getEnvString("GITHUB_CLIENT_SECRET", "")
	GitHubRedirectURL = getEnvString("GITHUB_REDIRECT_URL", fmt.Sprintf("http://localhost:%s/api/auth/github/callback", Port))
}

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

func GetGoogleOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     GoogleClientID,
		ClientSecret: GoogleClientSecret,
		RedirectURL:  GoogleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func GetGitHubOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     GitHubClientID,
		ClientSecret: GitHubClientSecret,
		RedirectURL:  GitHubRedirectURL,
		Scopes: []string{
			"user:email", // Access user email
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}
}
