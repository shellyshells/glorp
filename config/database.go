package config

import (
	"crypto/sha512"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDatabase() {
	var err error
	DB, err = sql.Open("sqlite3", "forum.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	createTables()
	runMigrations()
	seedDatabase()
	log.Println("Database initialized successfully")
}

func createTables() {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password_hash VARCHAR(128) NOT NULL,
		role VARCHAR(20) DEFAULT 'user',
		banned BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		last_login DATETIME,
		display_name VARCHAR(100),
		bio TEXT,
		location VARCHAR(100),
		website VARCHAR(255),
		avatar_url VARCHAR(255),
		show_email BOOLEAN DEFAULT FALSE,
		show_online BOOLEAN DEFAULT TRUE,
		allow_messages BOOLEAN DEFAULT TRUE,
		public_profile BOOLEAN DEFAULT TRUE,
		last_activity DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS tags (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(50) UNIQUE NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS threads (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title VARCHAR(200) NOT NULL,
		description TEXT NOT NULL,
		author_id INTEGER NOT NULL,
		status VARCHAR(20) DEFAULT 'open',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS thread_tags (
		thread_id INTEGER NOT NULL,
		tag_id INTEGER NOT NULL,
		PRIMARY KEY (thread_id, tag_id),
		FOREIGN KEY (thread_id) REFERENCES threads(id) ON DELETE CASCADE,
		FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		thread_id INTEGER NOT NULL,
		author_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (thread_id) REFERENCES threads(id) ON DELETE CASCADE,
		FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS votes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		message_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		vote_type INTEGER NOT NULL, -- 1 for like, -1 for dislike
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(message_id, user_id),
		FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS thread_votes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		thread_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		vote_type INTEGER NOT NULL, -- 1 for upvote, -1 for downvote
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(thread_id, user_id),
		FOREIGN KEY (thread_id) REFERENCES threads(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	`

	if _, err := DB.Exec(schema); err != nil {
		log.Fatal("Failed to create tables:", err)
	}
}

func runMigrations() {
	// Check if we need to add new columns to existing tables
	migrations := []string{
		"ALTER TABLE users ADD COLUMN display_name VARCHAR(100)",
		"ALTER TABLE users ADD COLUMN bio TEXT",
		"ALTER TABLE users ADD COLUMN location VARCHAR(100)",
		"ALTER TABLE users ADD COLUMN website VARCHAR(255)",
		"ALTER TABLE users ADD COLUMN avatar_url VARCHAR(255)",
		"ALTER TABLE users ADD COLUMN show_email BOOLEAN DEFAULT FALSE",
		"ALTER TABLE users ADD COLUMN show_online BOOLEAN DEFAULT TRUE",
		"ALTER TABLE users ADD COLUMN allow_messages BOOLEAN DEFAULT TRUE",
		"ALTER TABLE users ADD COLUMN public_profile BOOLEAN DEFAULT TRUE",
		"ALTER TABLE users ADD COLUMN last_activity DATETIME",
	}

	for _, migration := range migrations {
		_, err := DB.Exec(migration)
		if err != nil {
			// Column might already exist, continue
			log.Printf("Migration note: %v", err)
		}
	}

	// Update existing users with proper datetime values
	_, err := DB.Exec("UPDATE users SET display_name = username WHERE display_name IS NULL OR display_name = ''")
	if err != nil {
		log.Printf("Migration update error: %v", err)
	}

	// Set last_activity to created_at for existing users where it's NULL
	_, err = DB.Exec("UPDATE users SET last_activity = created_at WHERE last_activity IS NULL")
	if err != nil {
		log.Printf("Migration update error: %v", err)
	}
}

// Helper function to hash password using SHA512 (same as utils.HashPassword)
func hashPassword(password string) string {
	hash := sha512.Sum512([]byte(password))
	return fmt.Sprintf("%x", hash)
}

func seedDatabase() {
	// Check if admin user exists
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = 'admin'").Scan(&count)
	if err != nil {
		log.Printf("Failed to check admin user: %v", err)
		return
	}

	if count == 0 {
		// Create admin user with properly hashed password
		adminPassword := "AdminPassword123!"
		hashedPassword := hashPassword(adminPassword)

		log.Printf("🔑 Creating admin user with password: '%s'", adminPassword)
		log.Printf("🔒 Generated hash: '%s'", hashedPassword)

		_, err = DB.Exec(`INSERT INTO users (username, display_name, email, password_hash, role, created_at, last_activity) 
						 VALUES ('admin', 'admin', 'admin@forum.com', ?, 'admin', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`, hashedPassword)
		if err != nil {
			log.Printf("Failed to create admin user: %v", err)
		} else {
			log.Println("✅ Admin user created successfully")
		}

		// Insert some default tags
		tags := []string{"General", "Technology", "Programming", "Discussion", "Help", "News"}
		for _, tag := range tags {
			_, err = DB.Exec("INSERT OR IGNORE INTO tags (name) VALUES (?)", tag)
			if err != nil {
				log.Printf("Failed to insert tag %s: %v", tag, err)
			}
		}

		log.Println("Database seeded with admin user and default tags")
	} else {
		log.Println("✅ Admin user already exists")

		// Let's verify the existing admin user's password hash
		var existingHash string
		err = DB.QueryRow("SELECT password_hash FROM users WHERE username = 'admin'").Scan(&existingHash)
		if err == nil {
			adminPassword := "AdminPassword123!"
			correctHash := hashPassword(adminPassword)
			log.Printf("🔍 Existing admin hash: '%s'", existingHash)
			log.Printf("🔍 Correct hash should be: '%s'", correctHash)
			log.Printf("🔍 Hash matches: %v", existingHash == correctHash)

			// If the hash doesn't match, update it
			if existingHash != correctHash {
				log.Printf("🔧 Updating admin password hash...")
				_, err = DB.Exec("UPDATE users SET password_hash = ? WHERE username = 'admin'", correctHash)
				if err != nil {
					log.Printf("Failed to update admin password: %v", err)
				} else {
					log.Printf("✅ Admin password hash updated")
				}
			}
		}
	}

	// Debug: Check if users exist
	rows, err := DB.Query("SELECT username, email FROM users")
	if err != nil {
		log.Printf("Failed to query users: %v", err)
		return
	}
	defer rows.Close()

	log.Println("📋 Current users in database:")
	for rows.Next() {
		var username, email string
		if err := rows.Scan(&username, &email); err != nil {
			continue
		}
		log.Printf("  - %s (%s)", username, email)
	}
}
