package main

import (
	"log"
	"net/http"

	"goforum/config"
	"goforum/controllers"
	"goforum/middleware"

	"github.com/gorilla/mux"
)

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
	protected.HandleFunc("/user/{username}", controllers.UserByUsernameHandler).Methods("GET")

	// Admin routes
	admin := r.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.AuthMiddleware, middleware.AdminMiddleware)
	admin.HandleFunc("/dashboard", controllers.AdminDashboardHandler).Methods("GET")

	// API routes
	api := r.PathPrefix("/api").Subrouter()

	// Public API routes
	api.HandleFunc("/threads", controllers.GetThreadsHandler).Methods("GET")
	api.HandleFunc("/search", controllers.SearchHandler).Methods("GET")

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

	// Image upload API
	apiProtected.HandleFunc("/upload/image", controllers.UploadImageHandler).Methods("POST")
	apiProtected.HandleFunc("/upload/image/{filename}", controllers.DeleteImageHandler).Methods("DELETE")

	// Profile API
	apiProtected.HandleFunc("/profile/update", controllers.UpdateProfileHandler).Methods("POST")
	apiProtected.HandleFunc("/profile/avatar-style", controllers.UpdateAvatarHandler).Methods("POST")

	// Admin API routes
	apiAdmin := api.PathPrefix("/admin").Subrouter()
	apiAdmin.Use(middleware.AuthMiddleware, middleware.AdminMiddleware)
	apiAdmin.HandleFunc("/ban/{id:[0-9]+}", controllers.BanUserHandler).Methods("POST")
	apiAdmin.HandleFunc("/threads/{id:[0-9]+}/status", controllers.UpdateThreadStatusHandler).Methods("PUT")

	log.Println("🚀 GoForum server starting on :8080")
	log.Println("📱 Visit http://localhost:8080 to access the forum")
	log.Println("👤 Default admin: username=admin, password=AdminPassword123!")
	log.Fatal(http.ListenAndServe(":8080", r))
}
