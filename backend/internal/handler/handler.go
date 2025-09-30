package handler

import (
	"go-webapp/internal/store"

	"golang.org/x/oauth2"
)

type Handler struct {
	store             *store.Store
	googleOAuthConfig *oauth2.Config
	githubOAuthConfig *oauth2.Config
}

func NewHandler(store *store.Store, googleOAuthConfig *oauth2.Config, githubOAuthConfig *oauth2.Config) *Handler {
	return &Handler{store: store, googleOAuthConfig: googleOAuthConfig, githubOAuthConfig: githubOAuthConfig}
}
