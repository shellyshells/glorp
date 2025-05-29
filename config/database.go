package config

import (
	"database/sql"
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
		last_login DATETIME
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
	`

	if _, err := DB.Exec(schema); err != nil {
		log.Fatal("Failed to create tables:", err)
	}
}

func seedDatabase() {
	// Check if admin user exists
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = 'admin'").Scan(&count)
	if err != nil {
		log.Fatal("Failed to check admin user:", err)
	}

	if count == 0 {
		// Create admin user with hashed password
		hashedPassword := "ba3253876aed6bc22d4a6ff53d8406c6ad864195ed144ab5c87621b6c233b548baeae6956df346ec8c17f5ea10f35ee3cbc514797ed7ddd3145464e2a0bab413" // AdminPassword123!
		
		_, err = DB.Exec(`INSERT INTO users (username, email, password_hash, role) 
						 VALUES ('admin', 'admin@forum.com', ?, 'admin')`, hashedPassword)
		if err != nil {
			log.Fatal("Failed to create admin user:", err)
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
	}
}