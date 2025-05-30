package controllers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"goforum/middleware"
	"goforum/models"
	"goforum/utils"

	"github.com/gorilla/mux"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	currentUser := middleware.GetUserFromContext(r)
	if currentUser == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Get user to display (could be current user or someone else)
	vars := mux.Vars(r)
	username := vars["username"]

	var user *models.User
	var err error

	if username != "" {
		// Viewing someone else's profile
		user, err = models.GetUserByUsername(username)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
	} else {
		// Viewing own profile
		user, err = models.GetUserByID(currentUser.ID)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
	}

	// Update current user's activity
	models.UpdateUserActivity(currentUser.ID)

	// Get user's threads
	threadFilters := models.ThreadFilters{
		AuthorID: user.ID,
		Limit:    20,
		Page:     1,
	}
	userThreads, _, _ := models.GetThreads(threadFilters)

	// Get user's messages
	messageFilters := models.MessageFilters{
		UserID: user.ID,
		Limit:  20,
		Page:   1,
	}
	userMessages, _, _ := models.GetMessagesByUser(user.ID, messageFilters)

	// Get user stats
	threadCount, _ := models.GetThreadCountByUser(user.ID)
	messageCount, _ := models.GetMessageCountByUser(user.ID)

	// Calculate karma (this is a simplified version)
	postKarma := calculatePostKarma(userThreads)
	commentKarma := calculateCommentKarma(userMessages)

	tmpl := template.Must(template.New("").Funcs(TemplateFuncMap).ParseFiles("views/layouts/main.html", "views/user/profile.html"))
	data := map[string]interface{}{
		"Title":        "u/" + user.Username + " - Profile",
		"Page":         "profile",
		"User":         user,
		"CurrentUser":  currentUser,
		"UserThreads":  userThreads,
		"UserMessages": userMessages,
		"ThreadCount":  threadCount,
		"MessageCount": messageCount,
		"PostKarma":    postKarma,
		"CommentKarma": commentKarma,
	}

	if err := tmpl.ExecuteTemplate(w, "main.html", data); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func SettingsHandler(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Get fresh user data
	freshUser, err := models.GetUserByID(user.ID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	tmpl := template.Must(template.New("").Funcs(TemplateFuncMap).ParseFiles("views/layouts/main.html", "views/user/settings.html"))
	data := map[string]interface{}{
		"Title": "Settings - GoForum",
		"Page":  "settings",
		"User":  freshUser,
	}

	if err := tmpl.ExecuteTemplate(w, "main.html", data); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.UserProfile
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Sanitize input
	req.DisplayName = utils.SanitizeString(req.DisplayName)
	req.Bio = utils.SanitizeString(req.Bio)
	req.Location = utils.SanitizeString(req.Location)
	req.Website = utils.SanitizeString(req.Website)

	// Validate input
	if req.DisplayName == "" {
		req.DisplayName = user.Username
	}

	if len(req.Bio) > 500 {
		http.Error(w, "Bio must be less than 500 characters", http.StatusBadRequest)
		return
	}

	if len(req.Location) > 100 {
		http.Error(w, "Location must be less than 100 characters", http.StatusBadRequest)
		return
	}

	if req.Website != "" && !isValidURL(req.Website) {
		http.Error(w, "Invalid website URL", http.StatusBadRequest)
		return
	}

	// Update profile
	err := models.UpdateUserProfile(user.ID, req)
	if err != nil {
		log.Printf("Failed to update profile: %v", err)
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Profile updated successfully",
	})
}

func UpdateAvatarHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// For now, we'll just generate avatar URLs based on style
	// In a real app, you'd handle file uploads here
	var req struct {
		Style string `json:"style"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Generate avatar URL based on style
	avatarURL := generateAvatarURL(user.Username, req.Style)

	err := models.UpdateUserAvatar(user.ID, avatarURL)
	if err != nil {
		log.Printf("Failed to update avatar: %v", err)
		http.Error(w, "Failed to update avatar", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Avatar updated successfully",
	})
}

func UserByUsernameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	currentUser := middleware.GetUserFromContext(r)
	if currentUser == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Get user by username
	user, err := models.GetUserByUsername(username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if profile is public (unless it's the user themselves)
	if user.ID != currentUser.ID && !user.PublicProfile {
		http.Error(w, "This profile is private", http.StatusForbidden)
		return
	}

	// Update current user's activity
	models.UpdateUserActivity(currentUser.ID)

	// Get user's threads
	threadFilters := models.ThreadFilters{
		AuthorID: user.ID,
		Limit:    20,
		Page:     1,
	}
	userThreads, _, _ := models.GetThreads(threadFilters)

	// Get user's messages
	messageFilters := models.MessageFilters{
		UserID: user.ID,
		Limit:  20,
		Page:   1,
	}
	userMessages, _, _ := models.GetMessagesByUser(user.ID, messageFilters)

	// Get user stats
	threadCount, _ := models.GetThreadCountByUser(user.ID)
	messageCount, _ := models.GetMessageCountByUser(user.ID)

	// Calculate karma
	postKarma := calculatePostKarma(userThreads)
	commentKarma := calculateCommentKarma(userMessages)

	tmpl := template.Must(template.New("").Funcs(TemplateFuncMap).ParseFiles("views/layouts/main.html", "views/user/profile.html"))
	data := map[string]interface{}{
		"Title":        "u/" + user.Username + " - Profile",
		"Page":         "profile",
		"User":         user,
		"CurrentUser":  currentUser,
		"UserThreads":  userThreads,
		"UserMessages": userMessages,
		"ThreadCount":  threadCount,
		"MessageCount": messageCount,
		"PostKarma":    postKarma,
		"CommentKarma": commentKarma,
	}

	if err := tmpl.ExecuteTemplate(w, "main.html", data); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Helper functions
func calculatePostKarma(threads []models.Thread) int {
	karma := 0
	for _, thread := range threads {
		// In a real app, you'd get the actual vote scores for threads
		karma += thread.MessageCount // Simplified
	}
	return karma
}

func calculateCommentKarma(messages []models.Message) int {
	karma := 0
	for _, message := range messages {
		karma += message.Score
	}
	return karma
}

func isValidURL(url string) bool {
	// Simple URL validation
	return len(url) > 0 && (len(url) < 4 || url[:4] == "http")
}

func generateAvatarURL(username, style string) string {
	// Generate a style-based avatar URL
	// In a real app, you might use a service like Gravatar or generate actual images
	switch style {
	case "gradient1":
		return "/static/avatars/gradient1/" + username
	case "gradient2":
		return "/static/avatars/gradient2/" + username
	case "gradient3":
		return "/static/avatars/gradient3/" + username
	case "gradient4":
		return "/static/avatars/gradient4/" + username
	default:
		return ""
	}
}
