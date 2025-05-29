package controllers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"goforum/middleware"
	"goforum/models"
	"goforum/utils"

	"github.com/gorilla/mux"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	tagName := r.URL.Query().Get("tag")
	search := r.URL.Query().Get("search")
	sortBy := r.URL.Query().Get("sort")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, limit := utils.ParsePaginationParams(pageStr, limitStr)

	// Build filters
	filters := models.ThreadFilters{
		TagName: tagName,
		Search:  search,
		SortBy:  sortBy,
		Page:    page,
		Limit:   limit,
	}

	threads, total, err := models.GetThreads(filters)
	if err != nil {
		http.Error(w, "Failed to load threads", http.StatusInternalServerError)
		return
	}

	tags, _ := models.GetAllTags()
	pagination := utils.CalculatePagination(total, page, limit)

	// Check if user is authenticated
	user := middleware.GetUserFromContext(r)

	tmpl := template.Must(template.New("").Funcs(TemplateFuncMap).ParseFiles("views/layouts/main.html", "views/threads/index.html"))
	data := map[string]interface{}{
		"Title":      "GoForum - Home",
		"Page":       "home",
		"Threads":    threads,
		"Tags":       tags,
		"Pagination": pagination,
		"Filters":    filters,
		"User":       user,
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

	thread, err := models.GetThreadByID(threadID)
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

	user := middleware.GetUserFromContext(r)
	userID := 0
	if user != nil {
		userID = user.ID
	}

	messageFilters := models.MessageFilters{
		ThreadID: threadID,
		Page:     page,
		Limit:    limit,
		SortBy:   sortBy,
		UserID:   userID,
	}

	messages, totalMessages, err := models.GetMessages(messageFilters)
	if err != nil {
		http.Error(w, "Failed to load messages", http.StatusInternalServerError)
		return
	}

	pagination := utils.CalculatePagination(totalMessages, page, limit)

	tmpl := template.Must(template.New("").Funcs(TemplateFuncMap).ParseFiles("views/layouts/main.html", "views/threads/show.html"))
	data := map[string]interface{}{
		"Title":      thread.Title + " - GoForum",
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
	tags, _ := models.GetAllTags()
	user := middleware.GetUserFromContext(r)

	tmpl := template.Must(template.New("").Funcs(TemplateFuncMap).ParseFiles("views/layouts/main.html", "views/threads/create.html"))
	data := map[string]interface{}{
		"Title": "Create Thread - GoForum",
		"Page":  "create-thread",
		"Tags":  tags,
		"User":  user,
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
		"Title":  "Edit Thread - GoForum",
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
		Tags        []int  `json:"tags"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Sanitize input
	req.Title = utils.SanitizeString(req.Title)
	req.Description = utils.SanitizeString(req.Description)

	// Validate input
	if req.Title == "" || req.Description == "" {
		http.Error(w, "Title and description are required", http.StatusBadRequest)
		return
	}

	if len(req.Title) > 200 {
		http.Error(w, "Title must be less than 200 characters", http.StatusBadRequest)
		return
	}

	thread, err := models.CreateThread(req.Title, req.Description, user.ID, req.Tags)
	if err != nil {
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
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Sanitize input
	req.Title = utils.SanitizeString(req.Title)
	req.Description = utils.SanitizeString(req.Description)

	// Validate input
	if req.Title == "" || req.Description == "" {
		http.Error(w, "Title and description are required", http.StatusBadRequest)
		return
	}

	if len(req.Title) > 200 {
		http.Error(w, "Title must be less than 200 characters", http.StatusBadRequest)
		return
	}

	err = models.UpdateThread(threadID, req.Title, req.Description, req.Tags)
	if err != nil {
		http.Error(w, "Failed to update thread: "+err.Error(), http.StatusInternalServerError)
		return
	}

	thread, err := models.GetThreadByID(threadID)
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

	// Check permissions
	if !models.CanUserModifyThread(threadID, user.ID, user.Role) {
		http.Error(w, "Access denied", http.StatusForbidden)
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
