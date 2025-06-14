package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"glorp/config"
)

type User struct {
	ID            int        `json:"id"`
	Username      string     `json:"username"`
	DisplayName   string     `json:"display_name"`
	Email         string     `json:"email"`
	PasswordHash  string     `json:"-"`
	Role          string     `json:"role"`
	Banned        bool       `json:"banned"`
	Bio           string     `json:"bio"`
	Location      string     `json:"location"`
	Website       string     `json:"website"`
	AvatarURL     string     `json:"avatar_url"`
	AvatarStyle   string     `json:"avatar_style"`
	ShowEmail     bool       `json:"show_email"`
	ShowOnline    bool       `json:"show_online"`
	AllowMessages bool       `json:"allow_messages"`
	PublicProfile bool       `json:"public_profile"`
	CreatedAt     time.Time  `json:"created_at"`
	LastLogin     *time.Time `json:"last_login,omitempty"`
	LastActivity  time.Time  `json:"last_activity"`
}

type UserProfile struct {
	DisplayName   string `json:"display_name"`
	Bio           string `json:"bio"`
	Location      string `json:"location"`
	Website       string `json:"website"`
	AvatarStyle   string `json:"avatar_style"`
	ShowEmail     bool   `json:"show_email"`
	ShowOnline    bool   `json:"show_online"`
	AllowMessages bool   `json:"allow_messages"`
	PublicProfile bool   `json:"public_profile"`
}

func CreateUser(username, email, passwordHash string) (*User, error) {
	query := `INSERT INTO users (username, display_name, email, password_hash, avatar_style, created_at, last_activity) 
			  VALUES (?, ?, ?, ?, 'default', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`
	result, err := config.DB.Exec(query, username, username, email, passwordHash)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetUserByID(int(id))
}

func GetUserByID(id int) (*User, error) {
	user := &User{}

	log.Printf("🔍 Looking for user with ID: %d", id)

	query := `SELECT id, username, 
			         COALESCE(display_name, username) as display_name, 
			         email, password_hash, 
			         COALESCE(role, 'user') as role, 
			         COALESCE(banned, 0) as banned,
			         COALESCE(bio, '') as bio, 
			         COALESCE(location, '') as location, 
			         COALESCE(website, '') as website, 
			         COALESCE(avatar_url, '') as avatar_url,
			         COALESCE(avatar_style, 'default') as avatar_style,
			         COALESCE(show_email, 0) as show_email, 
			         COALESCE(show_online, 1) as show_online, 
			         COALESCE(allow_messages, 1) as allow_messages, 
			         COALESCE(public_profile, 1) as public_profile,
			         created_at, last_login, 
			         COALESCE(last_activity, created_at) as last_activity
			  FROM users WHERE id = ?`

	var lastLogin sql.NullTime
	var lastActivityStr sql.NullString
	err := config.DB.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.DisplayName, &user.Email, &user.PasswordHash,
		&user.Role, &user.Banned, &user.Bio, &user.Location, &user.Website, &user.AvatarURL,
		&user.AvatarStyle, &user.ShowEmail, &user.ShowOnline, &user.AllowMessages, &user.PublicProfile,
		&user.CreatedAt, &lastLogin, &lastActivityStr,
	)

	if err != nil {
		log.Printf("❌ Error scanning user with ID %d: %v", id, err)
		return nil, err
	}

	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}

	// Parse last_activity from string
	if lastActivityStr.Valid && lastActivityStr.String != "" {
		if parsedTime, err := time.Parse("2006-01-02 15:04:05", lastActivityStr.String); err == nil {
			user.LastActivity = parsedTime
		} else if parsedTime, err := time.Parse(time.RFC3339, lastActivityStr.String); err == nil {
			user.LastActivity = parsedTime
		} else {
			// Fallback to created_at if parsing fails
			user.LastActivity = user.CreatedAt
		}
	} else {
		user.LastActivity = user.CreatedAt
	}

	log.Printf("✅ Successfully loaded user: %s (ID: %d)", user.Username, user.ID)
	return user, nil
}

func GetUserByUsername(username string) (*User, error) {
	log.Printf("🔍 Looking for user with username: '%s'", username)

	query := `SELECT id, username, 
			         COALESCE(display_name, username) as display_name, 
			         email, password_hash, 
			         COALESCE(role, 'user') as role, 
			         COALESCE(banned, 0) as banned,
			         COALESCE(bio, '') as bio, 
			         COALESCE(location, '') as location, 
			         COALESCE(website, '') as website, 
			         COALESCE(avatar_url, '') as avatar_url,
			         COALESCE(avatar_style, 'default') as avatar_style,
			         COALESCE(show_email, 0) as show_email, 
			         COALESCE(show_online, 1) as show_online, 
			         COALESCE(allow_messages, 1) as allow_messages, 
			         COALESCE(public_profile, 1) as public_profile,
			         created_at, last_login, 
			         COALESCE(last_activity, created_at) as last_activity
			  FROM users WHERE username = ?`

	user := &User{}
	var lastLogin sql.NullTime
	var lastActivityStr sql.NullString
	err := config.DB.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.DisplayName, &user.Email, &user.PasswordHash,
		&user.Role, &user.Banned, &user.Bio, &user.Location, &user.Website, &user.AvatarURL,
		&user.AvatarStyle, &user.ShowEmail, &user.ShowOnline, &user.AllowMessages, &user.PublicProfile,
		&user.CreatedAt, &lastLogin, &lastActivityStr,
	)

	if err != nil {
		log.Printf("❌ Error scanning user with username '%s': %v", username, err)
		return nil, err
	}

	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}

	// Parse last_activity from string
	if lastActivityStr.Valid && lastActivityStr.String != "" {
		if parsedTime, err := time.Parse("2006-01-02 15:04:05", lastActivityStr.String); err == nil {
			user.LastActivity = parsedTime
		} else if parsedTime, err := time.Parse(time.RFC3339, lastActivityStr.String); err == nil {
			user.LastActivity = parsedTime
		} else {
			// Fallback to created_at if parsing fails
			user.LastActivity = user.CreatedAt
		}
	} else {
		user.LastActivity = user.CreatedAt
	}

	log.Printf("✅ Successfully loaded user: %s (ID: %d)", user.Username, user.ID)
	return user, nil
}

func GetUserByEmail(email string) (*User, error) {
	log.Printf("🔍 Looking for user with email: '%s'", email)

	query := `SELECT id, username, 
			         COALESCE(display_name, username) as display_name, 
			         email, password_hash, 
			         COALESCE(role, 'user') as role, 
			         COALESCE(banned, 0) as banned,
			         COALESCE(bio, '') as bio, 
			         COALESCE(location, '') as location, 
			         COALESCE(website, '') as website, 
			         COALESCE(avatar_url, '') as avatar_url,
			         COALESCE(avatar_style, 'default') as avatar_style,
			         COALESCE(show_email, 0) as show_email, 
			         COALESCE(show_online, 1) as show_online, 
			         COALESCE(allow_messages, 1) as allow_messages, 
			         COALESCE(public_profile, 1) as public_profile,
			         created_at, last_login, 
			         COALESCE(last_activity, created_at) as last_activity
			  FROM users WHERE email = ?`

	user := &User{}
	var lastLogin sql.NullTime
	var lastActivityStr sql.NullString
	err := config.DB.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.DisplayName, &user.Email, &user.PasswordHash,
		&user.Role, &user.Banned, &user.Bio, &user.Location, &user.Website, &user.AvatarURL,
		&user.AvatarStyle, &user.ShowEmail, &user.ShowOnline, &user.AllowMessages, &user.PublicProfile,
		&user.CreatedAt, &lastLogin, &lastActivityStr,
	)

	if err != nil {
		log.Printf("❌ Error scanning user with email '%s': %v", email, err)
		return nil, err
	}

	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}

	// Parse last_activity from string
	if lastActivityStr.Valid && lastActivityStr.String != "" {
		if parsedTime, err := time.Parse("2006-01-02 15:04:05", lastActivityStr.String); err == nil {
			user.LastActivity = parsedTime
		} else if parsedTime, err := time.Parse(time.RFC3339, lastActivityStr.String); err == nil {
			user.LastActivity = parsedTime
		} else {
			// Fallback to created_at if parsing fails
			user.LastActivity = user.CreatedAt
		}
	} else {
		user.LastActivity = user.CreatedAt
	}

	log.Printf("✅ Successfully loaded user: %s (ID: %d)", user.Username, user.ID)
	return user, nil
}

func GetUserByIdentifier(identifier string) (*User, error) {
	log.Printf("🔍 Looking for user with identifier: '%s'", identifier)

	// Try username first, then email
	user, err := GetUserByUsername(identifier)
	if err != nil {
		log.Printf("Username lookup failed, trying email...")
		user, err = GetUserByEmail(identifier)
	}
	return user, err
}

func UpdateUserProfile(userID int, profile UserProfile) error {
	query := `UPDATE users SET 
			  display_name = ?, bio = ?, location = ?, website = ?, avatar_style = ?,
			  show_email = ?, show_online = ?, allow_messages = ?, public_profile = ?
			  WHERE id = ?`

	_, err := config.DB.Exec(query,
		profile.DisplayName, profile.Bio, profile.Location, profile.Website, profile.AvatarStyle,
		profile.ShowEmail, profile.ShowOnline, profile.AllowMessages, profile.PublicProfile,
		userID)

	return err
}

func UpdateUserAvatar(userID int, avatarURL string) error {
	query := `UPDATE users SET avatar_url = ? WHERE id = ?`
	_, err := config.DB.Exec(query, avatarURL, userID)
	return err
}

func UpdateUserAvatarStyle(userID int, style string) error {
	query := `UPDATE users SET avatar_style = ? WHERE id = ?`
	_, err := config.DB.Exec(query, style, userID)
	return err
}

func UpdateUserLastLogin(userID int) error {
	query := `UPDATE users SET last_login = CURRENT_TIMESTAMP, last_activity = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := config.DB.Exec(query, userID)
	return err
}

func UpdateUserActivity(userID int) error {
	query := `UPDATE users SET last_activity = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := config.DB.Exec(query, userID)
	return err
}

func IsUserOnline(user *User) bool {
	if !user.ShowOnline {
		return false
	}

	// Consider user online if last activity was within 15 minutes
	threshold := time.Now().Add(-15 * time.Minute)
	return user.LastActivity.After(threshold)
}

func BanUser(userID int) error {
	query := `UPDATE users SET banned = TRUE WHERE id = ?`
	_, err := config.DB.Exec(query, userID)
	return err
}

func UnbanUser(userID int) error {
	query := "UPDATE users SET banned = 0 WHERE id = ?"
	_, err := config.DB.Exec(query, userID)
	return err
}

func DeleteUser(userID int) error {
	// Start a transaction for atomicity
	tx, err := config.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction for user deletion: %v", err)
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Recovered from panic during user deletion transaction: %v", r)
		}
	}()

	// Delete user's messages
	_, err = tx.Exec("DELETE FROM messages WHERE author_id = ?", userID)
	if err != nil {
		tx.Rollback()
		log.Printf("Error deleting messages for user %d: %v", userID, err)
		return fmt.Errorf("failed to delete user messages: %w", err)
	}

	// Delete user's threads
	_, err = tx.Exec("DELETE FROM threads WHERE author_id = ?", userID)
	if err != nil {
		tx.Rollback()
		log.Printf("Error deleting threads for user %d: %v", userID, err)
		return fmt.Errorf("failed to delete user threads: %w", err)
	}

	// Delete user's votes
	_, err = tx.Exec("DELETE FROM thread_votes WHERE user_id = ?", userID)
	if err != nil {
		tx.Rollback()
		log.Printf("Error deleting thread votes for user %d: %v", userID, err)
		return fmt.Errorf("failed to delete user thread votes: %w", err)
	}

	_, err = tx.Exec("DELETE FROM message_votes WHERE user_id = ?", userID)
	if err != nil {
		tx.Rollback()
		log.Printf("Error deleting message votes for user %d: %v", userID, err)
		return fmt.Errorf("failed to delete user message votes: %w", err)
	}

	// Delete user from communities
	_, err = tx.Exec("DELETE FROM community_members WHERE user_id = ?", userID)
	if err != nil {
		tx.Rollback()
		log.Printf("Error deleting community memberships for user %d: %v", userID, err)
		return fmt.Errorf("failed to delete user community memberships: %w", err)
	}

	// Finally, delete the user
	query := "DELETE FROM users WHERE id = ?"
	_, err = tx.Exec(query, userID)
	if err != nil {
		tx.Rollback()
		log.Printf("Error deleting user %d: %v", userID, err)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return tx.Commit()
}

func IsUsernameUnique(username string) bool {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE username = ?`
	config.DB.QueryRow(query, username).Scan(&count)
	return count == 0
}

func IsEmailUnique(email string) bool {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE email = ?`
	config.DB.QueryRow(query, email).Scan(&count)
	return count == 0
}

func ValidateUser(user *User) error {
	if user.Username == "" {
		return errors.New("username is required")
	}
	if user.Email == "" {
		return errors.New("email is required")
	}
	if len(user.Username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}
	if len(user.Username) > 50 {
		return errors.New("username must be less than 50 characters")
	}
	return nil
}

// GetUserInitial returns the first letter of the username for avatar display
func (u *User) GetUserInitial() string {
	if len(u.Username) > 0 {
		return string(u.Username[0])
	}
	return "U"
}

// GetAvatarStyle returns the avatar style or default
func (u *User) GetAvatarStyle() string {
	if u.AvatarStyle == "" {
		return "default"
	}
	return u.AvatarStyle
}

func GetUserCount() (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM users"
	err := config.DB.QueryRow(query).Scan(&count)
	if err != nil {
		log.Printf("Error getting user count: %v", err)
		return 0, err
	}
	return count, nil
}

func GetAllUsers() ([]User, error) {
	var users []User
	query := `SELECT id, username, 
			         COALESCE(display_name, username) as display_name, 
			         email, password_hash, 
			         COALESCE(role, 'user') as role, 
			         COALESCE(banned, 0) as banned,
			         COALESCE(bio, '') as bio, 
			         COALESCE(location, '') as location, 
			         COALESCE(website, '') as website, 
			         COALESCE(avatar_url, '') as avatar_url,
			         COALESCE(avatar_style, 'default') as avatar_style,
			         COALESCE(show_email, 0) as show_email, 
			         COALESCE(show_online, 1) as show_online, 
			         COALESCE(allow_messages, 1) as allow_messages, 
			         COALESCE(public_profile, 1) as public_profile,
			         created_at, last_login, 
			         COALESCE(last_activity, created_at) as last_activity
			  FROM users ORDER BY created_at DESC`

	rows, err := config.DB.Query(query)
	if err != nil {
		log.Printf("Error getting all users: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		var lastLogin sql.NullTime
		var lastActivityStr sql.NullString
		err := rows.Scan(
			&user.ID, &user.Username, &user.DisplayName, &user.Email, &user.PasswordHash,
			&user.Role, &user.Banned, &user.Bio, &user.Location, &user.Website, &user.AvatarURL,
			&user.AvatarStyle, &user.ShowEmail, &user.ShowOnline, &user.AllowMessages, &user.PublicProfile,
			&user.CreatedAt, &lastLogin, &lastActivityStr,
		)
		if err != nil {
			log.Printf("Error scanning user row: %v", err)
			continue
		}

		if lastLogin.Valid {
			user.LastLogin = &lastLogin.Time
		}

		// Parse last_activity from string
		if lastActivityStr.Valid && lastActivityStr.String != "" {
			if parsedTime, err := time.Parse("2006-01-02 15:04:05", lastActivityStr.String); err == nil {
				user.LastActivity = parsedTime
			} else if parsedTime, err := time.Parse(time.RFC3339, lastActivityStr.String); err == nil {
				user.LastActivity = parsedTime
			} else {
				// Fallback to created_at if parsing fails
				user.LastActivity = user.CreatedAt
			}
		} else {
			user.LastActivity = user.CreatedAt
		}
		users = append(users, user)
	}

	return users, nil
}
