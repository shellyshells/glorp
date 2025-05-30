package models

import (
	"database/sql"
	"errors"
	"time"

	"goforum/config"
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
	ShowEmail     bool   `json:"show_email"`
	ShowOnline    bool   `json:"show_online"`
	AllowMessages bool   `json:"allow_messages"`
	PublicProfile bool   `json:"public_profile"`
}

func CreateUser(username, email, passwordHash string) (*User, error) {
	query := `INSERT INTO users (username, display_name, email, password_hash, last_activity) VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)`
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
	query := `SELECT id, username, COALESCE(display_name, username), email, password_hash, role, banned, 
			         COALESCE(bio, ''), COALESCE(location, ''), COALESCE(website, ''), COALESCE(avatar_url, ''),
			         COALESCE(show_email, 0), COALESCE(show_online, 1), COALESCE(allow_messages, 1), COALESCE(public_profile, 1),
			         created_at, last_login, COALESCE(last_activity, created_at)
			  FROM users WHERE id = ?`

	var lastLogin sql.NullTime
	err := config.DB.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.DisplayName, &user.Email, &user.PasswordHash,
		&user.Role, &user.Banned, &user.Bio, &user.Location, &user.Website, &user.AvatarURL,
		&user.ShowEmail, &user.ShowOnline, &user.AllowMessages, &user.PublicProfile,
		&user.CreatedAt, &lastLogin, &user.LastActivity,
	)

	if err != nil {
		return nil, err
	}

	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}

	return user, nil
}

func GetUserByUsername(username string) (*User, error) {
	user := &User{}
	query := `SELECT id, username, COALESCE(display_name, username), email, password_hash, role, banned, 
			         COALESCE(bio, ''), COALESCE(location, ''), COALESCE(website, ''), COALESCE(avatar_url, ''),
			         COALESCE(show_email, 0), COALESCE(show_online, 1), COALESCE(allow_messages, 1), COALESCE(public_profile, 1),
			         created_at, last_login, COALESCE(last_activity, created_at)
			  FROM users WHERE username = ?`

	var lastLogin sql.NullTime
	err := config.DB.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.DisplayName, &user.Email, &user.PasswordHash,
		&user.Role, &user.Banned, &user.Bio, &user.Location, &user.Website, &user.AvatarURL,
		&user.ShowEmail, &user.ShowOnline, &user.AllowMessages, &user.PublicProfile,
		&user.CreatedAt, &lastLogin, &user.LastActivity,
	)

	if err != nil {
		return nil, err
	}

	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}

	return user, nil
}

func GetUserByEmail(email string) (*User, error) {
	user := &User{}
	query := `SELECT id, username, COALESCE(display_name, username), email, password_hash, role, banned, 
			         COALESCE(bio, ''), COALESCE(location, ''), COALESCE(website, ''), COALESCE(avatar_url, ''),
			         COALESCE(show_email, 0), COALESCE(show_online, 1), COALESCE(allow_messages, 1), COALESCE(public_profile, 1),
			         created_at, last_login, COALESCE(last_activity, created_at)
			  FROM users WHERE email = ?`

	var lastLogin sql.NullTime
	err := config.DB.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.DisplayName, &user.Email, &user.PasswordHash,
		&user.Role, &user.Banned, &user.Bio, &user.Location, &user.Website, &user.AvatarURL,
		&user.ShowEmail, &user.ShowOnline, &user.AllowMessages, &user.PublicProfile,
		&user.CreatedAt, &lastLogin, &user.LastActivity,
	)

	if err != nil {
		return nil, err
	}

	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}

	return user, nil
}

func GetUserByIdentifier(identifier string) (*User, error) {
	// Try username first, then email
	user, err := GetUserByUsername(identifier)
	if err != nil {
		user, err = GetUserByEmail(identifier)
	}
	return user, err
}

func UpdateUserProfile(userID int, profile UserProfile) error {
	query := `UPDATE users SET 
			  display_name = ?, bio = ?, location = ?, website = ?,
			  show_email = ?, show_online = ?, allow_messages = ?, public_profile = ?
			  WHERE id = ?`

	_, err := config.DB.Exec(query,
		profile.DisplayName, profile.Bio, profile.Location, profile.Website,
		profile.ShowEmail, profile.ShowOnline, profile.AllowMessages, profile.PublicProfile,
		userID)

	return err
}

func UpdateUserAvatar(userID int, avatarURL string) error {
	query := `UPDATE users SET avatar_url = ? WHERE id = ?`
	_, err := config.DB.Exec(query, avatarURL, userID)
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
	query := `UPDATE users SET banned = FALSE WHERE id = ?`
	_, err := config.DB.Exec(query, userID)
	return err
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
