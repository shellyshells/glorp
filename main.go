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

	// Static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// View routes (HTML)
	r.HandleFunc("/", controllers.HomeHandler).Methods("GET")
	r.HandleFunc("/register", controllers.RegisterViewHandler).Methods("GET")
	r.HandleFunc("/login", controllers.LoginViewHandler).Methods("GET")
	
	// Protected view routes
	protected := r.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/threads/create", controllers.CreateThreadViewHandler).Methods("GET")
	protected.HandleFunc("/threads/{id}/edit", controllers.EditThreadViewHandler).Methods("GET")
	
	// Admin routes
	admin := r.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.AuthMiddleware, middleware.AdminMiddleware)
	admin.HandleFunc("/dashboard", controllers.AdminDashboardHandler).Methods("GET")

	// Thread view (accessible to all)
	r.HandleFunc("/threads/{id}", controllers.ShowThreadHandler).Methods("GET")

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	
	// Auth API
	api.HandleFunc("/register", controllers.RegisterHandler).Methods("POST")
	api.HandleFunc("/login", controllers.LoginHandler).Methods("POST")
	api.HandleFunc("/logout", controllers.LogoutHandler).Methods("POST")

	// Protected API routes
	apiProtected := api.PathPrefix("").Subrouter()
	apiProtected.Use(middleware.AuthMiddleware)
	
	// Thread API
	apiProtected.HandleFunc("/threads", controllers.CreateThreadHandler).Methods("POST")
	apiProtected.HandleFunc("/threads/{id}", controllers.UpdateThreadHandler).Methods("PUT")
	apiProtected.HandleFunc("/threads/{id}", controllers.DeleteThreadHandler).Methods("DELETE")
	
	// Message API
	apiProtected.HandleFunc("/threads/{id}/messages", controllers.CreateMessageHandler).Methods("POST")
	apiProtected.HandleFunc("/messages/{id}", controllers.DeleteMessageHandler).Methods("DELETE")
	apiProtected.HandleFunc("/messages/{id}/vote", controllers.VoteMessageHandler).Methods("POST")

	// Public API routes
	api.HandleFunc("/threads", controllers.GetThreadsHandler).Methods("GET")
	api.HandleFunc("/search", controllers.SearchHandler).Methods("GET")

	// Admin API routes
	apiAdmin := api.PathPrefix("/admin").Subrouter()
	apiAdmin.Use(middleware.AuthMiddleware, middleware.AdminMiddleware)
	apiAdmin.HandleFunc("/ban/{id}", controllers.BanUserHandler).Methods("POST")
	apiAdmin.HandleFunc("/threads/{id}/status", controllers.UpdateThreadStatusHandler).Methods("PUT")

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}