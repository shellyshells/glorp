package controllers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"glorp/middleware"
	"glorp/models"

	"github.com/gorilla/mux"
)

func AdminDashboardHandler(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)

	// Get some statistics
	threadFilters := models.ThreadFilters{Limit: 0} // Get all threads
	threads, totalThreads, _ := models.GetThreads(threadFilters)

	totalMessages, _ := models.GetMessageCount()

	// Get recent threads
	recentThreadFilters := models.ThreadFilters{Limit: 10, Page: 1}
	recentThreads, _, _ := models.GetThreads(recentThreadFilters)

	tmpl := template.Must(template.New("").Funcs(TemplateFuncMap).ParseFiles("views/layouts/main.html", "views/admin/dashboard.html"))
	data := map[string]interface{}{
		"Title":         "Admin Dashboard - Glorp",
		"Page":          "admin-dashboard",
		"User":          user,
		"TotalThreads":  totalThreads,
		"TotalMessages": totalMessages,
		"RecentThreads": recentThreads,
		"Threads":       threads,
	}
	if err := tmpl.ExecuteTemplate(w, "main.html", data); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func BanUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil || user.Role != "admin" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Cannot ban yourself
	if userID == user.ID {
		http.Error(w, "Cannot ban yourself", http.StatusBadRequest)
		return
	}

	// Get the user to be banned
	targetUser, err := models.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Cannot ban other admins
	if targetUser.Role == "admin" {
		http.Error(w, "Cannot ban admin users", http.StatusForbidden)
		return
	}

	var req struct {
		Action string `json:"action"` // "ban" or "unban"
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var message string
	if req.Action == "ban" {
		err = models.BanUser(userID)
		message = "User banned successfully"
	} else if req.Action == "unban" {
		err = models.UnbanUser(userID)
		message = "User unbanned successfully"
	} else {
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Failed to update user status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": message,
	})
}

func UpdateThreadStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	threadID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid thread ID", http.StatusBadRequest)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil || user.Role != "admin" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Status string `json:"status"` // "open", "closed", "archived"
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate status
	validStatuses := map[string]bool{
		"open":     true,
		"closed":   true,
		"archived": true,
	}

	if !validStatuses[req.Status] {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	err = models.UpdateThreadStatus(threadID, req.Status)
	if err != nil {
		http.Error(w, "Failed to update thread status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Thread status updated successfully",
	})
}
