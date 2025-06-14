package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"glorp/middleware"
	"glorp/models"
	"glorp/utils"

	"github.com/gorilla/mux"
)

func CreateMessageHandler(w http.ResponseWriter, r *http.Request) {
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

	// Check if thread exists and is not closed/archived
	thread, err := models.GetThreadByID(threadID)
	if err != nil {
		http.Error(w, "Thread not found", http.StatusNotFound)
		return
	}

	if thread.Status == "closed" {
		http.Error(w, "Cannot post messages to closed thread", http.StatusForbidden)
		return
	}

	if thread.Status == "archived" {
		http.Error(w, "Thread not found", http.StatusNotFound)
		return
	}

	var req struct {
		Content  string `json:"content"`
		ParentID *int   `json:"parent_id"` // For nested replies
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Sanitize input
	req.Content = utils.SanitizeString(req.Content)

	// Validate input
	if req.Content == "" {
		http.Error(w, "Message content is required", http.StatusBadRequest)
		return
	}

	if len(req.Content) > 10000 {
		http.Error(w, "Message content is too long", http.StatusBadRequest)
		return
	}

	// If parent_id is provided, validate that the parent message exists and belongs to the same thread
	if req.ParentID != nil {
		parentMessage, err := models.GetMessageByID(*req.ParentID, 0)
		if err != nil {
			http.Error(w, "Parent message not found", http.StatusNotFound)
			return
		}
		if parentMessage.ThreadID != threadID {
			http.Error(w, "Parent message does not belong to this thread", http.StatusBadRequest)
			return
		}
	}

	message, err := models.CreateMessageWithParent(threadID, user.ID, req.Content, req.ParentID)
	if err != nil {
		http.Error(w, "Failed to create message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": message,
		"text":    "Message posted successfully",
	})
}

func VoteMessageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		VoteType int `json:"vote_type"` // 1 for like, -1 for dislike
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

	// Check if message exists
	message, err := models.GetMessageByID(messageID, user.ID)
	if err != nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	// Users cannot vote on their own messages
	if message.AuthorID == user.ID {
		http.Error(w, "Cannot vote on your own message", http.StatusForbidden)
		return
	}

	err = models.VoteMessage(messageID, user.ID, req.VoteType)
	if err != nil {
		http.Error(w, "Failed to vote: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get updated vote stats
	likes, dislikes, score, err := models.GetMessageVoteStats(messageID)
	if err != nil {
		http.Error(w, "Failed to get vote stats", http.StatusInternalServerError)
		return
	}

	// Get user's current vote
	userVote, _ := models.GetUserVote(messageID, user.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"likes":     likes,
		"dislikes":  dislikes,
		"score":     score,
		"user_vote": userVote,
		"message":   "Vote recorded successfully",
	})
}
