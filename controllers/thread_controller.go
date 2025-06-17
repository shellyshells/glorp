package controllers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"glorp/middleware"
	"glorp/models"
	"glorp/utils"

	"github.com/gorilla/mux"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	communityName := r.URL.Query().Get("community")
	search := r.URL.Query().Get("search")
	sortBy := r.URL.Query().Get("sort")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, limit := utils.ParsePaginationParams(pageStr, limitStr)

	// Build filters
	filters := models.ThreadFilters{
		Search: search,
		SortBy: sortBy,
		Page:   page,
		Limit:  limit,
	}

	// Add user ID to filters if user is authenticated
	user := middleware.GetUserFromContext(r)
	if user != nil {
		filters.UserID = user.ID
	}

	// If specific community is requested, get community ID
	var selectedCommunity *models.Community
	if communityName != "" {
		community, err := models.GetCommunityByName(communityName, filters.UserID)
		if err == nil {
			filters.CommunityID = community.ID
			selectedCommunity = community
		}
	}

	var threads []models.Thread
	var total int
	var err error

	if filters.CommunityID > 0 {
		// Get threads from specific community
		threads, total, err = models.GetThreadsByCommunity(filters)
	} else {
		// Get threads from all communities
		threads, total, err = models.GetThreads(filters)
	}

	if err != nil {
		http.Error(w, "Failed to load threads", http.StatusInternalServerError)
		return
	}

	// Get popular communities for sidebar
	communityFilters := models.CommunityFilters{
		Visibility: "public",
		SortBy:     "members",
		Limit:      10,
	}
	popularCommunities, _, _ := models.GetCommunities(communityFilters)

	// Get user's communities if authenticated
	var userCommunities []models.Community
	if user != nil {
		userCommunityFilters := models.CommunityFilters{
			UserID: user.ID,
			SortBy: "name",
			Limit:  20,
		}
		userCommunities, _, _ = models.GetCommunities(userCommunityFilters)
	}

	pagination := utils.CalculatePagination(total, page, limit)

	// Determine page title
	pageTitle := "Glorp - Home"
	if selectedCommunity != nil {
		pageTitle = "z/" + selectedCommunity.DisplayName + " - Glorp"
	}

	tmpl := template.Must(template.New("").Funcs(TemplateFuncMap).ParseFiles("views/layouts/main.html", "views/threads/index.html"))
	data := map[string]interface{}{
		"Title":              pageTitle,
		"Page":               "home",
		"Threads":            threads,
		"PopularCommunities": popularCommunities,
		"UserCommunities":    userCommunities,
		"SelectedCommunity":  selectedCommunity,
		"Pagination":         pagination,
		"Filters":            filters,
		"User":               user,
	}
	if err := tmpl.ExecuteTemplate(w, "main.html", data); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func ShowThreadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	threadID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid thread ID", http.StatusBadRequest)
		return
	}

	user := middleware.GetUserFromContext(r)
	userID := 0
	if user != nil {
		userID = user.ID
	}

	thread, err := models.GetThreadByIDWithUserAndCommunity(threadID, userID)
	if err != nil {
		http.Error(w, "Thread not found", http.StatusNotFound)
		return
	}

	// Check if thread is archived
	if thread.Status == "archived" {
		http.Error(w, "Thread not found", http.StatusNotFound)
		return
	}

	// Get messages
	sortBy := r.URL.Query().Get("sort")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, limit := utils.ParsePaginationParams(pageStr, limitStr)

	messageFilters := models.MessageFilters{
		ThreadID: threadID,
		Page:     page,
		Limit:    limit,
		SortBy:   sortBy,
		UserID:   userID,
	}

	messages, totalMessages, err := models.GetMessages(messageFilters)
	if err != nil {
		log.Printf("Error loading messages for thread %d: %v", threadID, err)
		http.Error(w, "Failed to load messages", http.StatusInternalServerError)
		return
	}

	// Get user roles for the community if thread belongs to a community
	if thread.CommunityID != nil {
		if user != nil {
			// Get user's role in the community
			userRole, err := models.GetUserCommunityRole(*thread.CommunityID, user.ID)
			if err == nil {
				user.Role = userRole
			}
		}
	}

	pagination := utils.CalculatePagination(totalMessages, page, limit)

	tmpl := template.Must(template.New("").Funcs(TemplateFuncMap).ParseFiles("views/layouts/main.html", "views/threads/show.html"))
	data := map[string]interface{}{
		"Title":      thread.Title + " - Glorp",
		"Page":       "thread",
		"Thread":     thread,
		"Messages":   messages,
		"Pagination": pagination,
		"User":       user,
		"SortBy":     sortBy,
	}
	if err := tmpl.ExecuteTemplate(w, "main.html", data); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func CreateThreadViewHandler(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Get all communities the user can post to
	var availableCommunities []models.Community

	// Get user's communities (communities they're a member of)
	userCommunityFilters := models.CommunityFilters{
		UserID: user.ID,
		SortBy: "name",
		Limit:  100,
	}
	userCommunities, _, err := models.GetCommunities(userCommunityFilters)
	if err == nil {
		availableCommunities = append(availableCommunities, userCommunities...)
	}

	// Also get public communities that allow open posting
	publicCommunityFilters := models.CommunityFilters{
		Visibility: "public",
		SortBy:     "name",
		Limit:      100,
	}
	publicCommunities, _, err := models.GetCommunities(publicCommunityFilters)
	if err == nil {
		// Add public communities that aren't already in the list
		communityMap := make(map[int]bool)
		for _, community := range availableCommunities {
			communityMap[community.ID] = true
		}

		for _, community := range publicCommunities {
			if !communityMap[community.ID] {
				availableCommunities = append(availableCommunities, community)
			}
		}
	}

	// Get all tags for backward compatibility
	tags, _ := models.GetAllTags()

	// Check if specific community was requested via URL parameter
	requestedCommunity := r.URL.Query().Get("community")

	tmpl := template.Must(template.New("").Funcs(TemplateFuncMap).ParseFiles("views/layouts/main.html", "views/threads/create.html"))
	data := map[string]interface{}{
		"Title":              "Create Thread - Glorp",
		"Page":               "create-thread",
		"Communities":        availableCommunities,
		"Tags":               tags,
		"RequestedCommunity": requestedCommunity,
		"User":               user,
	}
	if err := tmpl.ExecuteTemplate(w, "main.html", data); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func EditThreadViewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	threadID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid thread ID", http.StatusBadRequest)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	thread, err := models.GetThreadByID(threadID)
	if err != nil {
		http.Error(w, "Thread not found", http.StatusNotFound)
		return
	}

	// Check permissions
	if !models.CanUserModifyThread(threadID, user.ID, user.Role) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	tags, _ := models.GetAllTags()

	tmpl := template.Must(template.New("").Funcs(TemplateFuncMap).ParseFiles("views/layouts/main.html", "views/threads/edit.html"))
	data := map[string]interface{}{
		"Title":  "Edit Thread - Glorp",
		"Page":   "edit-thread",
		"Thread": thread,
		"Tags":   tags,
		"User":   user,
	}
	if err := tmpl.ExecuteTemplate(w, "main.html", data); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func GetThreadsHandler(w http.ResponseWriter, r *http.Request) {
	tagName := r.URL.Query().Get("tag")
	search := r.URL.Query().Get("search")
	sortBy := r.URL.Query().Get("sort")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, limit := utils.ParsePaginationParams(pageStr, limitStr)

	filters := models.ThreadFilters{
		TagName: tagName,
		Search:  search,
		SortBy:  sortBy,
		Page:    page,
		Limit:   limit,
	}

	// Add user ID to filters if user is authenticated
	user := middleware.GetUserFromContext(r)
	if user != nil {
		filters.UserID = user.ID
	}

	threads, total, err := models.GetThreads(filters)
	if err != nil {
		http.Error(w, "Failed to load threads", http.StatusInternalServerError)
		return
	}

	pagination := utils.CalculatePagination(total, page, limit)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"threads":    threads,
		"pagination": pagination,
	})
}

func CreateThreadHandler(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		CommunityID int    `json:"community_id"` // Community instead of tags
		PostType    string `json:"post_type"`    // text, link, image
		ImageURL    string `json:"image_url"`    // For image posts
		LinkURL     string `json:"link_url"`     // For link posts
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Sanitize input
	req.Title = utils.SanitizeString(req.Title)
	req.Description = utils.SanitizeString(req.Description)
	req.ImageURL = utils.SanitizeString(req.ImageURL)
	req.LinkURL = utils.SanitizeString(req.LinkURL)

	// Validate input
	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	if len(req.Title) > 200 {
		http.Error(w, "Title must be less than 200 characters", http.StatusBadRequest)
		return
	}

	// Set default post type
	if req.PostType == "" {
		req.PostType = "text"
	}

	// Validate post type
	validPostTypes := map[string]bool{
		"text":  true,
		"link":  true,
		"image": true,
	}

	if !validPostTypes[req.PostType] {
		http.Error(w, "Invalid post type", http.StatusBadRequest)
		return
	}

	// Validate community selection
	if req.CommunityID <= 0 {
		http.Error(w, "Community selection is required", http.StatusBadRequest)
		return
	}

	// Check if user can post in this community
	canPost, err := models.CanUserPostInCommunity(req.CommunityID, user.ID)
	if err != nil {
		http.Error(w, "Failed to verify community permissions", http.StatusInternalServerError)
		return
	}

	if !canPost {
		http.Error(w, "You don't have permission to post in this community", http.StatusForbidden)
		return
	}

	// For text posts, description is optional but for link posts we need either description or URL
	if req.PostType == "link" && req.LinkURL == "" && req.Description == "" {
		http.Error(w, "Link posts require either a URL or description", http.StatusBadRequest)
		return
	}

	// For image posts, we need either an image URL or description
	if req.PostType == "image" && req.ImageURL == "" && req.Description == "" {
		http.Error(w, "Image posts require either an image or description", http.StatusBadRequest)
		return
	}

	// Create the thread in the specified community
	thread, err := models.CreateThreadInCommunity(req.Title, req.Description, user.ID, req.CommunityID, req.PostType, req.ImageURL, req.LinkURL)
	if err != nil {
		log.Printf("Error creating thread: %v", err)
		http.Error(w, "Failed to create thread: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"thread":  thread,
		"message": "Thread created successfully",
	})
}

func UpdateThreadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	threadID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid thread ID", http.StatusBadRequest)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check permissions
	if !models.CanUserModifyThread(threadID, user.ID, user.Role) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Tags        []int  `json:"tags"`
		ImageURL    string `json:"image_url"`
		LinkURL     string `json:"link_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Sanitize input
	req.Title = utils.SanitizeString(req.Title)
	req.Description = utils.SanitizeString(req.Description)
	req.ImageURL = utils.SanitizeString(req.ImageURL)
	req.LinkURL = utils.SanitizeString(req.LinkURL)

	// Validate input
	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	if len(req.Title) > 200 {
		http.Error(w, "Title must be less than 200 characters", http.StatusBadRequest)
		return
	}

	err = models.UpdateThreadWithMedia(threadID, req.Title, req.Description, req.ImageURL, req.LinkURL, req.Tags)
	if err != nil {
		http.Error(w, "Failed to update thread: "+err.Error(), http.StatusInternalServerError)
		return
	}

	thread, err := models.GetThreadByIDWithUser(threadID, user.ID)
	if err != nil {
		http.Error(w, "Failed to get updated thread", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"thread":  thread,
		"message": "Thread updated successfully",
	})
}

func DeleteThreadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	threadID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid thread ID", http.StatusBadRequest)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user is authorized to delete the thread (either admin or author)
	if !models.CanUserModifyThread(threadID, user.ID, user.Role) {
		http.Error(w, "Unauthorized to delete this thread", http.StatusForbidden)
		return
	}

	err = models.DeleteThread(threadID)
	if err != nil {
		http.Error(w, "Failed to delete thread: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Thread deleted successfully",
	})
}

func VoteThreadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	threadID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid thread ID", http.StatusBadRequest)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		VoteType int `json:"vote_type"` // 1 for upvote, -1 for downvote
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate vote type
	if req.VoteType != 1 && req.VoteType != -1 {
		http.Error(w, "Invalid vote type", http.StatusBadRequest)
		return
	}

	// Check if thread exists
	_, err = models.GetThreadByID(threadID)
	if err != nil {
		http.Error(w, "Thread not found", http.StatusNotFound)
		return
	}

	// Users can vote on their own threads
	err = models.VoteThread(threadID, user.ID, req.VoteType)
	if err != nil {
		http.Error(w, "Failed to vote: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get updated vote stats
	upvotes, downvotes, score, err := models.GetThreadVoteStats(threadID)
	if err != nil {
		http.Error(w, "Failed to get vote stats", http.StatusInternalServerError)
		return
	}

	// Get user's current vote
	userVote, _ := models.GetUserThreadVote(threadID, user.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"upvotes":   upvotes,
		"downvotes": downvotes,
		"score":     score,
		"user_vote": userVote,
		"message":   "Vote recorded successfully",
	})
}
