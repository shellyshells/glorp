package middleware

import (
	"context"
	"net/http"

	"goforum/utils"
)

type contextKey string

const UserContextKey contextKey = "user"

type UserContext struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for JWT token in cookie
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// Validate JWT token
		claims, err := utils.ValidateJWT(cookie.Value)
		if err != nil {
			// Clear invalid cookie
			http.SetCookie(w, &http.Cookie{
				Name:   "auth_token",
				Value:  "",
				MaxAge: -1,
				Path:   "/",
			})
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// Add user info to context
		userCtx := UserContext{
			ID:       claims.UserID,
			Username: claims.Username,
			Role:     claims.Role,
		}

		ctx := context.WithValue(r.Context(), UserContextKey, userCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(r *http.Request) *UserContext {
	if user, ok := r.Context().Value(UserContextKey).(UserContext); ok {
		return &user
	}
	return nil
}

func OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for JWT token in cookie
		cookie, err := r.Cookie("auth_token")
		if err == nil {
			// Validate JWT token
			claims, err := utils.ValidateJWT(cookie.Value)
			if err == nil {
				// Add user info to context
				userCtx := UserContext{
					ID:       claims.UserID,
					Username: claims.Username,
					Role:     claims.Role,
				}
				ctx := context.WithValue(r.Context(), UserContextKey, userCtx)
				r = r.WithContext(ctx)
			}
		}

		next.ServeHTTP(w, r)
	})
}