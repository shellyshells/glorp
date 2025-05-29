package models

import (
	"database/sql"
	"errors"
	"time"

	"goforum/config"
)

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	Banned       bool      `json:"banned"`
	CreatedAt    time.Time `json:"created_at"`
	LastLogin    *time.Time `json:"last_login,omitempty"`
}

func CreateUser(username, email, passwordHash string) (*User, error) {
	query := `INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)`
	result, err := config.DB.Exec(query, username, email, passwordHash)
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
	query := `SELECT id, username, email, password_hash, role, banned, created_at, last_login 
			  FROM users WHERE id = ?`
	
	var lastLogin sql.NullTime
	err := config.DB.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Role, &user.Banned, &user.CreatedAt, &lastLogin,
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
	query := `SELECT id, username, email, password_hash, role, banned, created_at, last_login 
			  FROM users WHERE username = ?`
	
	var lastLogin sql.NullTime
	err := config.DB.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Role, &user.Banned, &user.CreatedAt, &lastLogin,
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
	query := `SELECT id, username, email, password_hash, role, banned, created_at, last_login 
			  FROM users WHERE email = ?`
	
	var lastLogin sql.NullTime
	err := config.DB.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Role, &user.Banned, &user.CreatedAt, &lastLogin,
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

func UpdateUserLastLogin(userID int) error {
	query := `UPDATE users SET last_login = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := config.DB.Exec(query, userID)
	return err
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