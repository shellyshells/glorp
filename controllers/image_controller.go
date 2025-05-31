package controllers

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"goforum/middleware"

	"github.com/gorilla/mux"
)

// UploadImageHandler handles image uploads
func UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse multipart form (max 10MB)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "File too large (max 10MB)", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "No image file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file type
	if !isValidImageType(header.Filename) {
		http.Error(w, "Invalid file type. Only JPG, PNG, GIF, and WEBP are allowed", http.StatusBadRequest)
		return
	}

	// Validate file size (max 10MB)
	if header.Size > 10<<20 {
		http.Error(w, "File too large (max 10MB)", http.StatusBadRequest)
		return
	}

	// Create uploads directory if it doesn't exist
	uploadsDir := "static/uploads"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	// Generate unique filename
	filename := generateUniqueFilename(header.Filename, user.ID)
	filepath := filepath.Join(uploadsDir, filename)

	// Create the file
	dst, err := os.Create(filepath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the destination
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Return the file URL
	imageURL := fmt.Sprintf("/static/uploads/%s", filename)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"success": true, "url": "%s", "filename": "%s"}`, imageURL, filename)
}

// isValidImageType checks if the file has a valid image extension
func isValidImageType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}

	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}
	return false
}

// generateUniqueFilename creates a unique filename using timestamp and hash
func generateUniqueFilename(originalFilename string, userID int) string {
	ext := filepath.Ext(originalFilename)
	timestamp := time.Now().Unix()

	// Create hash from original filename and user ID
	hash := md5.Sum([]byte(fmt.Sprintf("%s_%d_%d", originalFilename, userID, timestamp)))
	hashStr := fmt.Sprintf("%x", hash)[:8]

	return fmt.Sprintf("%d_%s_%d%s", userID, hashStr, timestamp, ext)
}

// DeleteImageHandler handles image deletion (optional, for cleanup)
func DeleteImageHandler(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	filename := vars["filename"]

	// Security check: only allow deletion of files uploaded by the user
	if !strings.HasPrefix(filename, fmt.Sprintf("%d_", user.ID)) && user.Role != "admin" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	filepath := filepath.Join("static/uploads", filename)

	// Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Delete the file
	if err := os.Remove(filepath); err != nil {
		http.Error(w, "Failed to delete file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"success": true, "message": "Image deleted successfully"}`)
}
