package middleware

import (
	"context"
	"net/http"

	"glorp/models"
	"glorp/utils"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for JWT token in cookie
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			// Check if this is an API request
			if isAPIRequest(r) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			// For web requests, redirect to login
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

			if isAPIRequest(r) {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// Get full user from database to check if banned and get latest info
		user, err := models.GetUserByID(claims.UserID)
		if err != nil {
			if isAPIRequest(r) {
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// Check if user is banned
		if user.Banned {
			// Clear cookie and redirect/error
			http.SetCookie(w, &http.Cookie{
				Name:   "auth_token",
				Value:  "",
				MaxAge: -1,
				Path:   "/",
			})

			if isAPIRequest(r) {
				http.Error(w, "Account has been banned", http.StatusForbidden)
				return
			}
			http.Redirect(w, r, "/login?error=banned", http.StatusFound)
			return
		}

		// Add FULL user object to context (not UserContext)
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for JWT token in cookie
		cookie, err := r.Cookie("auth_token")
		if err == nil {
			// Validate JWT token
			claims, err := utils.ValidateJWT(cookie.Value)
			if err == nil {
				// Get full user from database
				user, err := models.GetUserByID(claims.UserID)
				if err == nil && !user.Banned {
					// Add FULL user object to context
					ctx := context.WithValue(r.Context(), UserContextKey, user)
					r = r.WithContext(ctx)
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

func GetUserFromContext(r *http.Request) *models.User {
	if user, ok := r.Context().Value(UserContextKey).(*models.User); ok {
		return user
	}
	return nil
}

func isAPIRequest(r *http.Request) bool {
	return r.Header.Get("Content-Type") == "application/json" ||
		r.Header.Get("Accept") == "application/json" ||
		(len(r.URL.Path) >= 4 && r.URL.Path[:4] == "/api")
}
