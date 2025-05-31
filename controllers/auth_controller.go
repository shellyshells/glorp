package controllers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"goforum/models"
	"goforum/utils"
)

func RegisterViewHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("").Funcs(TemplateFuncMap).ParseFiles("views/layouts/main.html", "views/auth/register.html"))
	data := map[string]interface{}{
		"Title": "Register - GoForum",
		"Page":  "register",
	}
	if err := tmpl.ExecuteTemplate(w, "main.html", data); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func LoginViewHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("").Funcs(TemplateFuncMap).ParseFiles("views/layouts/main.html", "views/auth/login.html"))
	data := map[string]interface{}{
		"Title": "Login - GoForum",
		"Page":  "login",
	}
	if err := tmpl.ExecuteTemplate(w, "main.html", data); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("JSON decode error: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Sanitize input
	req.Username = utils.SanitizeString(req.Username)
	req.Email = utils.SanitizeString(req.Email)

	// Validate input
	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Validate email format
	if err := utils.ValidateEmail(req.Email); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate password strength
	if err := utils.ValidatePassword(req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if username is unique
	if !models.IsUsernameUnique(req.Username) {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}

	// Check if email is unique
	if !models.IsEmailUnique(req.Email) {
		http.Error(w, "Email already exists", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword := utils.HashPassword(req.Password)

	// Create user
	user, err := models.CreateUser(req.Username, req.Email, hashedPassword)
	if err != nil {
		log.Printf("Create user error: %v", err)
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Username, user.Role)
	if err != nil {
		log.Printf("JWT generation error: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   86400, // 24 hours
	})

	// Update last login
	models.UpdateUserLastLogin(user.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"user":    user,
		"message": "Registration successful",
	})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Identifier string `json:"identifier"` // username or email
		Password   string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("JSON decode error: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Sanitize input
	req.Identifier = utils.SanitizeString(req.Identifier)

	// Validate input
	if req.Identifier == "" || req.Password == "" {
		log.Printf("Login validation failed: identifier=%s, password_length=%d", req.Identifier, len(req.Password))
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Get user by identifier (username or email)
	user, err := models.GetUserByIdentifier(req.Identifier)
	if err != nil {
		log.Printf("User lookup error for identifier '%s': %v", req.Identifier, err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Check if user is banned
	if user.Banned {
		log.Printf("Banned user login attempt: %s", user.Username)
		http.Error(w, "Account has been banned", http.StatusForbidden)
		return
	}

	// DEBUG: Log password verification details
	log.Printf("🔒 DEBUG Password verification for user '%s':", user.Username)
	log.Printf("   - Input password: '%s'", req.Password)
	log.Printf("   - Stored hash: '%s'", user.PasswordHash)

	// Generate hash for the input password to compare
	inputPasswordHash := utils.HashPassword(req.Password)
	log.Printf("   - Generated hash: '%s'", inputPasswordHash)
	log.Printf("   - Hashes match: %v", inputPasswordHash == user.PasswordHash)

	// Verify password
	if !utils.VerifyPassword(req.Password, user.PasswordHash) {
		log.Printf("Password verification failed for user: %s", user.Username)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Username, user.Role)
	if err != nil {
		log.Printf("JWT generation error for user %s: %v", user.Username, err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   86400, // 24 hours
	})

	// Update last login
	err = models.UpdateUserLastLogin(user.ID)
	if err != nil {
		log.Printf("Failed to update last login for user %s: %v", user.Username, err)
		// Don't fail the login for this
	}

	log.Printf("Successful login for user: %s", user.Username)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"user":    user,
		"message": "Login successful",
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear auth cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "auth_token",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Logout successful",
	})
}
