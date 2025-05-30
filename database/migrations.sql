-- Add profile and online status columns to users table
ALTER TABLE users ADD COLUMN display_name VARCHAR(100);
ALTER TABLE users ADD COLUMN bio TEXT;
ALTER TABLE users ADD COLUMN location VARCHAR(100);
ALTER TABLE users ADD COLUMN website VARCHAR(255);
ALTER TABLE users ADD COLUMN avatar_url VARCHAR(255);
ALTER TABLE users ADD COLUMN show_email BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN show_online BOOLEAN DEFAULT TRUE;
ALTER TABLE users ADD COLUMN allow_messages BOOLEAN DEFAULT TRUE;
ALTER TABLE users ADD COLUMN public_profile BOOLEAN DEFAULT TRUE;
ALTER TABLE users ADD COLUMN last_activity DATETIME DEFAULT CURRENT_TIMESTAMP;

-- Add parent_id to messages for nested replies
ALTER TABLE messages ADD COLUMN parent_id INTEGER;
ALTER TABLE messages ADD FOREIGN KEY (parent_id) REFERENCES messages(id) ON DELETE CASCADE;

-- Add thread voting table
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

-- Update existing users with display_name
UPDATE users SET display_name = username WHERE display_name IS NULL;