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
		avatar_style VARCHAR(50) DEFAULT 'default',
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
		post_type VARCHAR(20) DEFAULT 'text',
		image_url VARCHAR(500),
		link_url VARCHAR(500),
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
		parent_id INTEGER,
		author_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (thread_id) REFERENCES threads(id) ON DELETE CASCADE,
		FOREIGN KEY (parent_id) REFERENCES messages(id) ON DELETE CASCADE,
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

	CREATE TABLE IF NOT EXISTS uploaded_files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		filename VARCHAR(255) NOT NULL,
		original_name VARCHAR(255) NOT NULL,
		file_size INTEGER NOT NULL,
		mime_type VARCHAR(100) NOT NULL,
		user_id INTEGER NOT NULL,
		thread_id INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (thread_id) REFERENCES threads(id) ON DELETE SET NULL
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
		"ALTER TABLE users ADD COLUMN avatar_style VARCHAR(50) DEFAULT 'default'",
		"ALTER TABLE users ADD COLUMN show_email BOOLEAN DEFAULT FALSE",
		"ALTER TABLE users ADD COLUMN show_online BOOLEAN DEFAULT TRUE",
		"ALTER TABLE users ADD COLUMN allow_messages BOOLEAN DEFAULT TRUE",
		"ALTER TABLE users ADD COLUMN public_profile BOOLEAN DEFAULT TRUE",
		"ALTER TABLE users ADD COLUMN last_activity DATETIME",
		"ALTER TABLE threads ADD COLUMN post_type VARCHAR(20) DEFAULT 'text'",
		"ALTER TABLE threads ADD COLUMN image_url VARCHAR(500)",
		"ALTER TABLE threads ADD COLUMN link_url VARCHAR(500)",
		"ALTER TABLE messages ADD COLUMN parent_id INTEGER",
	}

	for _, migration := range migrations {
		_, err := DB.Exec(migration)
		if err != nil {
			// Column might already exist, continue
			log.Printf("Migration note: %v", err)
		}
	}

	// Update existing data
	updateExistingData()

	// Create indexes for better performance
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_threads_post_type ON threads(post_type)",
		"CREATE INDEX IF NOT EXISTS idx_threads_author ON threads(author_id)",
		"CREATE INDEX IF NOT EXISTS idx_threads_status ON threads(status)",
		"CREATE INDEX IF NOT EXISTS idx_threads_created ON threads(created_at)",
		"CREATE INDEX IF NOT EXISTS idx_messages_thread ON messages(thread_id)",
		"CREATE INDEX IF NOT EXISTS idx_messages_parent ON messages(parent_id)",
		"CREATE INDEX IF NOT EXISTS idx_messages_author ON messages(author_id)",
		"CREATE INDEX IF NOT EXISTS idx_messages_created ON messages(created_at)",
		"CREATE INDEX IF NOT EXISTS idx_votes_message ON votes(message_id)",
		"CREATE INDEX IF NOT EXISTS idx_votes_user ON votes(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_thread_votes_thread ON thread_votes(thread_id)",
		"CREATE INDEX IF NOT EXISTS idx_thread_votes_user ON thread_votes(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_thread_tags_thread ON thread_tags(thread_id)",
		"CREATE INDEX IF NOT EXISTS idx_thread_tags_tag ON thread_tags(tag_id)",
		"CREATE INDEX IF NOT EXISTS idx_uploaded_files_user ON uploaded_files(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_uploaded_files_thread ON uploaded_files(thread_id)",
	}

	for _, index := range indexes {
		_, err := DB.Exec(index)
		if err != nil {
			log.Printf("Index creation note: %v", err)
		}
	}
}

func updateExistingData() {
	// Update existing users with proper values
	_, err := DB.Exec("UPDATE users SET display_name = username WHERE display_name IS NULL OR display_name = ''")
	if err != nil {
		log.Printf("Migration update error: %v", err)
	}

	// Set last_activity to created_at for existing users where it's NULL
	_, err = DB.Exec("UPDATE users SET last_activity = created_at WHERE last_activity IS NULL")
	if err != nil {
		log.Printf("Migration update error: %v", err)
	}

	// Update existing threads to have 'text' post_type
	_, err = DB.Exec("UPDATE threads SET post_type = 'text' WHERE post_type IS NULL OR post_type = ''")
	if err != nil {
		log.Printf("Migration update error: %v", err)
	}

	// Set default avatar style for existing users
	_, err = DB.Exec("UPDATE users SET avatar_style = 'default' WHERE avatar_style IS NULL OR avatar_style = ''")
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

		_, err = DB.Exec(`INSERT INTO users (username, display_name, email, password_hash, role, avatar_style, created_at, last_activity) 
						 VALUES ('admin', 'admin', 'admin@forum.com', ?, 'admin', 'admin', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`, hashedPassword)
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
