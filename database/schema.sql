-- Glorp Database Schema
-- This file documents the database structure used by the application

-- Users table - stores user account information
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(128) NOT NULL,
    role VARCHAR(20) DEFAULT 'user',
    banned BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_login DATETIME
);

-- Tags table - stores available tags/categories
CREATE TABLE tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(50) UNIQUE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Threads table - stores discussion threads
CREATE TABLE threads (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    author_id INTEGER NOT NULL,
    status VARCHAR(20) DEFAULT 'open', -- 'open', 'closed', 'archived'
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Thread_tags junction table - many-to-many relationship between threads and tags
CREATE TABLE thread_tags (
    thread_id INTEGER NOT NULL,
    tag_id INTEGER NOT NULL,
    PRIMARY KEY (thread_id, tag_id),
    FOREIGN KEY (thread_id) REFERENCES threads(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

-- Messages table - stores messages within threads
CREATE TABLE messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    thread_id INTEGER NOT NULL,
    author_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (thread_id) REFERENCES threads(id) ON DELETE CASCADE,
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Votes table - stores likes/dislikes for messages
CREATE TABLE votes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    message_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    vote_type INTEGER NOT NULL, -- 1 for like, -1 for dislike
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(message_id, user_id), -- Ensure one vote per user per message
    FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Indexes for better performance
CREATE INDEX idx_threads_author ON threads(author_id);
CREATE INDEX idx_threads_status ON threads(status);
CREATE INDEX idx_threads_created ON threads(created_at);
CREATE INDEX idx_messages_thread ON messages(thread_id);
CREATE INDEX idx_messages_author ON messages(author_id);
CREATE INDEX idx_messages_created ON messages(created_at);
CREATE INDEX idx_votes_message ON votes(message_id);
CREATE INDEX idx_votes_user ON votes(user_id);
CREATE INDEX idx_thread_tags_thread ON thread_tags(thread_id);
CREATE INDEX idx_thread_tags_tag ON thread_tags(tag_id);