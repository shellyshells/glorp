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

-- Add new tables for community/subreddit functionality

-- Communities table (replaces the simple tags system)
CREATE TABLE IF NOT EXISTS communities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(50) UNIQUE NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    description TEXT,
    creator_id INTEGER NOT NULL,
    visibility VARCHAR(20) DEFAULT 'public', -- 'public', 'private', 'restricted'
    join_approval VARCHAR(20) DEFAULT 'open', -- 'open', 'approval_required', 'invite_only'
    member_count INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Community memberships
CREATE TABLE IF NOT EXISTS community_memberships (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    community_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    role VARCHAR(20) DEFAULT 'member', -- 'member', 'moderator', 'admin', 'creator'
    status VARCHAR(20) DEFAULT 'active', -- 'active', 'pending', 'banned'
    joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(community_id, user_id),
    FOREIGN KEY (community_id) REFERENCES communities(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Community join requests
CREATE TABLE IF NOT EXISTS community_join_requests (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    community_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    message TEXT,
    status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'approved', 'rejected'
    requested_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    processed_at DATETIME,
    processed_by INTEGER,
    UNIQUE(community_id, user_id),
    FOREIGN KEY (community_id) REFERENCES communities(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (processed_by) REFERENCES users(id) ON DELETE SET NULL
);

-- Community rules
CREATE TABLE IF NOT EXISTS community_rules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    community_id INTEGER NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    rule_order INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (community_id) REFERENCES communities(id) ON DELETE CASCADE
);

-- Update threads table to use communities instead of tags
ALTER TABLE threads ADD COLUMN community_id INTEGER;
ALTER TABLE threads ADD FOREIGN KEY (community_id) REFERENCES communities(id) ON DELETE SET NULL;

-- Migration script to convert existing tags to communities
INSERT INTO communities (name, display_name, description, creator_id, member_count)
SELECT 
    LOWER(REPLACE(name, ' ', '')), 
    name, 
    'Auto-created community from tag: ' || name,
    1, -- Admin user ID
    (SELECT COUNT(DISTINCT t.author_id) 
     FROM threads t 
     JOIN thread_tags tt ON t.id = tt.thread_id 
     WHERE tt.tag_id = tags.id)
FROM tags;

-- Update existing threads with community IDs
UPDATE threads 
SET community_id = (
    SELECT c.id 
    FROM communities c 
    JOIN tags t ON LOWER(REPLACE(t.name, ' ', '')) = c.name
    JOIN thread_tags tt ON t.id = tt.tag_id
    WHERE tt.thread_id = threads.id
    LIMIT 1
);

-- Add all existing thread authors as members of their communities
INSERT OR IGNORE INTO community_memberships (community_id, user_id, role)
SELECT DISTINCT 
    threads.community_id,
    threads.author_id,
    'member'
FROM threads 
WHERE threads.community_id IS NOT NULL;

-- Add admin as creator/moderator of all auto-created communities
INSERT OR IGNORE INTO community_memberships (community_id, user_id, role)
SELECT id, creator_id, 'creator'
FROM communities;

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_communities_name ON communities(name);
CREATE INDEX IF NOT EXISTS idx_communities_visibility ON communities(visibility);
CREATE INDEX IF NOT EXISTS idx_community_memberships_community ON community_memberships(community_id);
CREATE INDEX IF NOT EXISTS idx_community_memberships_user ON community_memberships(user_id);
CREATE INDEX IF NOT EXISTS idx_community_memberships_role ON community_memberships(role);
CREATE INDEX IF NOT EXISTS idx_community_join_requests_community ON community_join_requests(community_id);
CREATE INDEX IF NOT EXISTS idx_community_join_requests_user ON community_join_requests(user_id);
CREATE INDEX IF NOT EXISTS idx_community_join_requests_status ON community_join_requests(status);
CREATE INDEX IF NOT EXISTS idx_threads_community ON threads(community_id);