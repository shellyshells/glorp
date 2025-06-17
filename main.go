package main

import (
	"fmt"
	"log"
	"net/http"

	"glorp/config"
	"glorp/controllers"
	"glorp/middleware"
	"glorp/utils"

	"github.com/gorilla/mux"
)

// Error handling middleware
func errorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				utils.ShowErrorPage(w, r, http.StatusInternalServerError,
					"Something went wrong. Please try again later.",
					fmt.Sprintf("%v", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// 404 handler
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	utils.ShowErrorPage(w, r, http.StatusNotFound,
		"The page you're looking for doesn't exist or has been moved.",
		"")
}

func main() {
	// Initialize database
	config.InitDatabase()

	// Create router
	r := mux.NewRouter()

	// Static files (including uploaded images)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Apply optional auth middleware to home and thread view routes
	homeRoutes := r.PathPrefix("").Subrouter()
	homeRoutes.Use(middleware.OptionalAuthMiddleware)
	homeRoutes.HandleFunc("/", controllers.HomeHandler).Methods("GET")
	homeRoutes.HandleFunc("/threads/{id:[0-9]+}", controllers.ShowThreadHandler).Methods("GET")

	// Community routes (with optional auth)
	homeRoutes.HandleFunc("/communities", controllers.CommunityListHandler).Methods("GET")
	homeRoutes.HandleFunc("/z/{name}", controllers.CommunityViewHandler).Methods("GET")

	// Public auth routes (no middleware)
	r.HandleFunc("/register", controllers.RegisterViewHandler).Methods("GET")
	r.HandleFunc("/login", controllers.LoginViewHandler).Methods("GET")

	// Protected view routes
	protected := r.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/threads/create", controllers.CreateThreadViewHandler).Methods("GET")
	protected.HandleFunc("/threads/{id:[0-9]+}/edit", controllers.EditThreadViewHandler).Methods("GET")
	protected.HandleFunc("/profile", controllers.ProfileHandler).Methods("GET")
	protected.HandleFunc("/settings", controllers.SettingsHandler).Methods("GET")

	// Protected community routes
	protected.HandleFunc("/communities/create", controllers.CreateCommunityViewHandler).Methods("GET")
	protected.HandleFunc("/z/{name}/manage", controllers.CommunityManageHandler).Methods("GET")

	// User profile routes - these need to be accessible to view other users' profiles
	userRoutes := r.PathPrefix("/user").Subrouter()
	userRoutes.Use(middleware.AuthMiddleware)
	userRoutes.HandleFunc("/{username}", controllers.UserByUsernameHandler).Methods("GET")

	// Admin routes
	admin := r.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.AuthMiddleware, middleware.AdminMiddleware)
	admin.HandleFunc("/dashboard", controllers.AdminDashboardHandler).Methods("GET")

	// API routes
	api := r.PathPrefix("/api").Subrouter()

	// Public API routes
	api.HandleFunc("/threads", controllers.GetThreadsHandler).Methods("GET")
	api.HandleFunc("/search", controllers.SearchHandler).Methods("GET")
	api.HandleFunc("/communities", controllers.GetCommunitiesHandler).Methods("GET")

	// Auth API routes (no middleware required)
	api.HandleFunc("/register", controllers.RegisterHandler).Methods("POST")
	api.HandleFunc("/login", controllers.LoginHandler).Methods("POST")
	api.HandleFunc("/logout", controllers.LogoutHandler).Methods("POST")

	// Protected API routes
	apiProtected := api.PathPrefix("").Subrouter()
	apiProtected.Use(middleware.AuthMiddleware)

	// Thread API
	apiProtected.HandleFunc("/threads", controllers.CreateThreadHandler).Methods("POST")
	apiProtected.HandleFunc("/threads/{id:[0-9]+}", controllers.UpdateThreadHandler).Methods("PUT")
	apiProtected.HandleFunc("/threads/{id:[0-9]+}", controllers.DeleteThreadHandler).Methods("DELETE")
	apiProtected.HandleFunc("/threads/{id:[0-9]+}/vote", controllers.VoteThreadHandler).Methods("POST")

	// Message API
	apiProtected.HandleFunc("/threads/{id:[0-9]+}/messages", controllers.CreateMessageHandler).Methods("POST")
	apiProtected.HandleFunc("/messages/{id:[0-9]+}", controllers.DeleteMessageHandler).Methods("DELETE")
	apiProtected.HandleFunc("/messages/{id:[0-9]+}/vote", controllers.VoteMessageHandler).Methods("POST")

	// Community API
	apiProtected.HandleFunc("/communities", controllers.CreateCommunityHandler).Methods("POST")
	apiProtected.HandleFunc("/communities/{id:[0-9]+}", controllers.UpdateCommunityHandler).Methods("PUT")
	apiProtected.HandleFunc("/communities/{id:[0-9]+}/join", controllers.JoinCommunityHandler).Methods("POST")
	apiProtected.HandleFunc("/communities/{id:[0-9]+}/leave", controllers.LeaveCommunityHandler).Methods("POST")
	apiProtected.HandleFunc("/communities/join-requests/{id:[0-9]+}", controllers.ProcessJoinRequestHandler).Methods("POST")
	apiProtected.HandleFunc("/communities/{communityId:[0-9]+}/moderators/{userId:[0-9]+}", controllers.ManageModeratorHandler).Methods("POST")

	// Image upload API
	apiProtected.HandleFunc("/upload/image", controllers.UploadImageHandler).Methods("POST")
	apiProtected.HandleFunc("/upload/image/{filename}", controllers.DeleteImageHandler).Methods("DELETE")

	// Profile API
	apiProtected.HandleFunc("/profile/update", controllers.UpdateProfileHandler).Methods("POST")
	apiProtected.HandleFunc("/profile/avatar", controllers.UpdateAvatarHandler).Methods("POST")
	apiProtected.HandleFunc("/profile/avatar-style", controllers.UpdateAvatarHandler).Methods("POST")

	// Admin API routes
	apiAdmin := api.PathPrefix("/admin").Subrouter()
	apiAdmin.Use(middleware.AuthMiddleware, middleware.AdminMiddleware)
	apiAdmin.HandleFunc("/ban/{id:[0-9]+}", controllers.BanUserHandler).Methods("POST")
	apiAdmin.HandleFunc("/threads/{id:[0-9]+}/status", controllers.UpdateThreadStatusHandler).Methods("PUT")
	apiAdmin.HandleFunc("/users/{id:[0-9]+}", controllers.DeleteUserHandler).Methods("DELETE")
	apiAdmin.HandleFunc("/messages/{id:[0-9]+}", controllers.DeleteMessageHandler).Methods("DELETE")
	apiAdmin.HandleFunc("/messages/{id:[0-9]+}", controllers.EditMessageHandler).Methods("PUT")
	apiAdmin.HandleFunc("/communities/{id:[0-9]+}", controllers.DeleteCommunityHandler).Methods("DELETE")
	apiAdmin.HandleFunc("/communities/{id:[0-9]+}", controllers.UpdateCommunityHandler).Methods("PUT")

	// Add error handling middleware
	r.Use(errorHandler)

	// Set up 404 handler
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	log.Println("🚀 Glorp server starting on :8080")
	log.Println("📱 Visit http://localhost:8080 to access the forum")
	log.Println("👤 Default admin: username=admin, password=AdminPassword123!")
	log.Println("🏘️  Community system enabled!")
	log.Fatal(http.ListenAndServe(":8080", r))
}
