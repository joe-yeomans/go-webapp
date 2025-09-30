package util

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// generateSecureCode generates a cryptographically secure 6-digit code
func GenerateSecureCode() (string, error) {
	// Generate a random number between 0 and 999999
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	// Format as 6-digit string with leading zeros
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func SetSessionCookie(w http.ResponseWriter, sessionToken string) {
	// Use secure cookies only in production (HTTPS)
	frontendUrl := GetEnvString("FRONTEND_URL", "http://localhost:3000")
	cookieDomain := GetEnvString("COOKIE_DOMAIN", "localhost")
	isSecure := strings.HasPrefix(frontendUrl, "https://")

	sessionTimeMinutes := GetEnvInt("SESSION_TIME_MINUTES", 10)
	sessionTime := time.Duration(sessionTimeMinutes) * time.Minute

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

func ClearSessionCookie(w http.ResponseWriter) {
	// Use secure cookies only in production (HTTPS)
	frontendUrl := GetEnvString("FRONTEND_URL", "http://localhost:3000")
	cookieDomain := GetEnvString("COOKIE_DOMAIN", "localhost")
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
