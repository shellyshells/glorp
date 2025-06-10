package controllers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"goforum/config"
	"goforum/middleware"
	"goforum/models"
	"goforum/utils"

	"github.com/gorilla/mux"
)

// Community list page
func CommunityListHandler(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	search := r.URL.Query().Get("search")
	visibility := r.URL.Query().Get("visibility")
	sortBy := r.URL.Query().Get("sort")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, limit := utils.ParsePaginationParams(pageStr, limitStr)

	// Build filters
	filters := models.CommunityFilters{
		Search:     search,
		Visibility: visibility,
		SortBy:     sortBy,
		Page:       page,
		Limit:      limit,
	}

	user := middleware.GetUserFromContext(r)
	communities, total, err := models.GetCommunities(filters)
	if err != nil {
		http.Error(w, "Failed to load communities", http.StatusInternalServerError)
		return
	}

	pagination := utils.CalculatePagination(total, page, limit)

	tmpl := template.Must(template.New("").Funcs(TemplateFuncMap).ParseFiles("views/layouts/main.html", "views/communities/index.html"))
	data := map[string]interface{}{
		"Title":       "Communities - GoForum",
		"Page":        "communities",
		"Communities": communities,
		"Pagination":  pagination,
		"Filters":     filters,
		"User":        user,
	}
	if err := tmpl.ExecuteTemplate(w, "main.html", data); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Community view page
func CommunityViewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	communityName := vars["name"]

	user := middleware.GetUserFromContext(r)
	userID := 0
	if user != nil {
		userID = user.ID
	}

	community, err := models.GetCommunityByName(communityName, userID)
	if err != nil {
		http.Error(w, "Community not found", http.StatusNotFound)
		return
	}

	// Check if user can view this community
	if community.Visibility == "private" && community.UserRole == "" {
		http.Error(w, "This community is private", http.StatusForbidden)
		return
	}

	// Get community threads
	threadFilters := models.ThreadFilters{
		CommunityID: community.ID,
		Page:        1,
		Limit:       25,
		UserID:      userID,
	}

	threads, totalThreads, err := models.GetThreadsByCommunity(threadFilters)
	if err != nil {
		http.Error(w, "Failed to load threads", http.StatusInternalServerError)
		return
	}

	// Get community moderators
	moderators, _ := models.GetCommunityModerators(community.ID)

	// Get pending join requests if user can manage community
	var pendingRequests []models.CommunityJoinRequest
	if models.CanManageCommunity(community.ID, userID) {
		pendingRequests, _ = models.GetPendingJoinRequests(community.ID)
	}

	tmpl := template.Must(template.New("").Funcs(TemplateFuncMap).ParseFiles("views/layouts/main.html", "views/communities/show.html"))
	data := map[string]interface{}{
		"Title":           "r/" + community.DisplayName + " - GoForum",
		"Page":            "community",
		"Community":       community,
		"Threads":         threads,
		"TotalThreads":    totalThreads,
		"Moderators":      moderators,
		"PendingRequests": pendingRequests,
		"User":            user,
	}
	if err := tmpl.ExecuteTemplate(w, "main.html", data); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Create community page
func CreateCommunityViewHandler(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	tmpl := template.Must(template.New("").Funcs(TemplateFuncMap).ParseFiles("views/layouts/main.html", "views/communities/create.html"))
	data := map[string]interface{}{
		"Title": "Create Community - GoForum",
		"Page":  "create-community",
		"User":  user,
	}
	if err := tmpl.ExecuteTemplate(w, "main.html", data); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Community management page
func CommunityManageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	communityName := vars["name"]

	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	community, err := models.GetCommunityByName(communityName, user.ID)
	if err != nil {
		http.Error(w, "Community not found", http.StatusNotFound)
		return
	}

	// Check if user can manage this community
	if !models.CanManageCommunity(community.ID, user.ID) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Get community data for management
	moderators, _ := models.GetCommunityModerators(community.ID)
	pendingRequests, _ := models.GetPendingJoinRequests(community.ID)

	tmpl := template.Must(template.New("").Funcs(TemplateFuncMap).ParseFiles("views/layouts/main.html", "views/communities/manage.html"))
	data := map[string]interface{}{
		"Title":           "Manage r/" + community.DisplayName + " - GoForum",
		"Page":            "community-manage",
		"Community":       community,
		"Moderators":      moderators,
		"PendingRequests": pendingRequests,
		"User":            user,
	}
	if err := tmpl.ExecuteTemplate(w, "main.html", data); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// API Handlers

// Get communities API
func GetCommunitiesHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	visibility := r.URL.Query().Get("visibility")
	sortBy := r.URL.Query().Get("sort")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, limit := utils.ParsePaginationParams(pageStr, limitStr)

	filters := models.CommunityFilters{
		Search:     search,
		Visibility: visibility,
		SortBy:     sortBy,
		Page:       page,
		Limit:      limit,
	}

	communities, total, err := models.GetCommunities(filters)
	if err != nil {
		http.Error(w, "Failed to load communities", http.StatusInternalServerError)
		return
	}

	pagination := utils.CalculatePagination(total, page, limit)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"communities": communities,
		"pagination":  pagination,
	})
}

// Create community API
func CreateCommunityHandler(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Name         string `json:"name"`
		DisplayName  string `json:"display_name"`
		Description  string `json:"description"`
		Visibility   string `json:"visibility"`    // public, private, restricted
		JoinApproval string `json:"join_approval"` // open, approval_required, invite_only
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Sanitize input
	req.Name = utils.SanitizeString(req.Name)
	req.DisplayName = utils.SanitizeString(req.DisplayName)
	req.Description = utils.SanitizeString(req.Description)

	// Validate input
	if req.Name == "" {
		http.Error(w, "Community name is required", http.StatusBadRequest)
		return
	}

	if req.DisplayName == "" {
		req.DisplayName = req.Name
	}

	// Set defaults
	if req.Visibility == "" {
		req.Visibility = "public"
	}
	if req.JoinApproval == "" {
		req.JoinApproval = "open"
	}

	// Validate visibility and join approval
	validVisibilities := map[string]bool{"public": true, "private": true, "restricted": true}
	validJoinApprovals := map[string]bool{"open": true, "approval_required": true, "invite_only": true}

	if !validVisibilities[req.Visibility] {
		http.Error(w, "Invalid visibility setting", http.StatusBadRequest)
		return
	}

	if !validJoinApprovals[req.JoinApproval] {
		http.Error(w, "Invalid join approval setting", http.StatusBadRequest)
		return
	}

	community, err := models.CreateCommunity(req.Name, req.DisplayName, req.Description, user.ID, req.Visibility, req.JoinApproval)
	if err != nil {
		log.Printf("Error creating community: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"community": community,
		"message":   "Community created successfully",
	})
}

// Join community API
func JoinCommunityHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	communityID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid community ID", http.StatusBadRequest)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Message string `json:"message"` // For approval requests
	}

	json.NewDecoder(r.Body).Decode(&req)

	err = models.JoinCommunity(communityID, user.ID, req.Message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Successfully joined community",
	})
}

// Leave community API
func LeaveCommunityHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	communityID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid community ID", http.StatusBadRequest)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = models.LeaveCommunity(communityID, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Successfully left community",
	})
}

// Process join request API
func ProcessJoinRequestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid request ID", http.StatusBadRequest)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Approved bool `json:"approved"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = models.ProcessJoinRequest(requestID, user.ID, req.Approved)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	action := "rejected"
	if req.Approved {
		action = "approved"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Join request " + action + " successfully",
	})
}

// Update community settings API
func UpdateCommunityHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	communityID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid community ID", http.StatusBadRequest)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check permissions
	if !models.CanManageCommunity(communityID, user.ID) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	var req struct {
		DisplayName  string `json:"display_name"`
		Description  string `json:"description"`
		Visibility   string `json:"visibility"`
		JoinApproval string `json:"join_approval"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Sanitize input
	req.DisplayName = utils.SanitizeString(req.DisplayName)
	req.Description = utils.SanitizeString(req.Description)

	// Validate input
	if req.DisplayName == "" {
		http.Error(w, "Display name is required", http.StatusBadRequest)
		return
	}

	// Validate visibility and join approval
	validVisibilities := map[string]bool{"public": true, "private": true, "restricted": true}
	validJoinApprovals := map[string]bool{"open": true, "approval_required": true, "invite_only": true}

	if req.Visibility != "" && !validVisibilities[req.Visibility] {
		http.Error(w, "Invalid visibility setting", http.StatusBadRequest)
		return
	}

	if req.JoinApproval != "" && !validJoinApprovals[req.JoinApproval] {
		http.Error(w, "Invalid join approval setting", http.StatusBadRequest)
		return
	}

	// Update community
	query := `UPDATE communities SET display_name = ?, description = ?, updated_at = CURRENT_TIMESTAMP`
	args := []interface{}{req.DisplayName, req.Description}

	if req.Visibility != "" {
		query += ", visibility = ?"
		args = append(args, req.Visibility)
	}

	if req.JoinApproval != "" {
		query += ", join_approval = ?"
		args = append(args, req.JoinApproval)
	}

	query += " WHERE id = ?"
	args = append(args, communityID)

	_, err = config.DB.Exec(query, args...)
	if err != nil {
		http.Error(w, "Failed to update community", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Community updated successfully",
	})
}

// Add/remove moderator API
func ManageModeratorHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	communityID, err := strconv.Atoi(vars["communityId"])
	if err != nil {
		http.Error(w, "Invalid community ID", http.StatusBadRequest)
		return
	}

	targetUserID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check permissions (only creator can manage moderators)
	var creatorID int
	err = config.DB.QueryRow("SELECT creator_id FROM communities WHERE id = ?", communityID).Scan(&creatorID)
	if err != nil || creatorID != user.ID {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	var req struct {
		Action string `json:"action"` // "add" or "remove"
		Role   string `json:"role"`   // "moderator" or "admin"
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Role == "" {
		req.Role = "moderator"
	}

	validRoles := map[string]bool{"moderator": true, "admin": true}
	if !validRoles[req.Role] {
		http.Error(w, "Invalid role", http.StatusBadRequest)
		return
	}

	if req.Action == "add" {
		// Update user's role
		_, err = config.DB.Exec(`
			UPDATE community_memberships 
			SET role = ? 
			WHERE community_id = ? AND user_id = ? AND status = 'active'`,
			req.Role, communityID, targetUserID)
	} else if req.Action == "remove" {
		// Demote to regular member
		_, err = config.DB.Exec(`
			UPDATE community_memberships 
			SET role = 'member' 
			WHERE community_id = ? AND user_id = ? AND status = 'active'`,
			communityID, targetUserID)
	} else {
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Failed to update moderator status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Moderator status updated successfully",
	})
}
