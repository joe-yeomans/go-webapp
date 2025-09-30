package handler

import (
	"go-webapp/internal/store"
	"go-webapp/internal/util"
	"log"
	"net/http"
)

func (h *Handler) HandleLogoutRequest(w http.ResponseWriter, r *http.Request) {
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
		h.store.DeleteSession(cookie.Value)
		log.Printf("Session deleted for user %s", user.Email)
	}

	util.ClearSessionCookie(w)

	// Get frontend URL
	frontendUrl := util.GetEnvString("FRONTEND_URL", "http://localhost:3000")

	log.Printf("User %s successfully logged out", user.Email)

	// Redirect to frontend home page
	http.Redirect(w, r, frontendUrl, http.StatusSeeOther)
}
