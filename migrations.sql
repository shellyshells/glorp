-- Migration script to add community system to existing glorp database
-- Run this script to upgrade from tag-based to community-based system

BEGIN TRANSACTION;

-- 1. Create new community tables
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

CREATE TABLE IF NOT EXISTS community_rules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    community_id INTEGER NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    rule_order INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (community_id) REFERENCES communities(id) ON DELETE CASCADE
);

-- 2. Add community_id to threads table
ALTER TABLE threads ADD COLUMN community_id INTEGER;
ALTER TABLE threads ADD CONSTRAINT fk_threads_community 
    FOREIGN KEY (community_id) REFERENCES communities(id) ON DELETE SET NULL;

-- 3. Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_communities_name ON communities(name);
CREATE INDEX IF NOT EXISTS idx_communities_visibility ON communities(visibility);
CREATE INDEX IF NOT EXISTS idx_community_memberships_community ON community_memberships(community_id);
CREATE INDEX IF NOT EXISTS idx_community_memberships_user ON community_memberships(user_id);
CREATE INDEX IF NOT EXISTS idx_community_memberships_role ON community_memberships(role);
CREATE INDEX IF NOT EXISTS idx_community_join_requests_community ON community_join_requests(community_id);
CREATE INDEX IF NOT EXISTS idx_community_join_requests_user ON community_join_requests(user_id);
CREATE INDEX IF NOT EXISTS idx_community_join_requests_status ON community_join_requests(status);
CREATE INDEX IF NOT EXISTS idx_threads_community ON threads(community_id);

-- 4. Migrate existing tags to communities (only if tags table exists and has data)
INSERT OR IGNORE INTO communities (name, display_name, description, creator_id, member_count)
SELECT 
    LOWER(REPLACE(REPLACE(name, ' ', ''), '-', '_')) as name,
    name as display_name,
    'Auto-created community from tag: ' || name as description,
    1 as creator_id, -- Admin user ID
    COALESCE((
        SELECT COUNT(DISTINCT t.author_id) 
        FROM threads t 
        JOIN thread_tags tt ON t.id = tt.thread_id 
        WHERE tt.tag_id = tags.id
    ), 1) as member_count
FROM tags 
WHERE EXISTS (SELECT 1 FROM sqlite_master WHERE type='table' AND name='tags');

-- 5. Update existing threads with community IDs (only if migration is needed)
UPDATE threads 
SET community_id = (
    SELECT c.id 
    FROM communities c 
    JOIN tags t ON LOWER(REPLACE(REPLACE(t.name, ' ', ''), '-', '_')) = c.name
    JOIN thread_tags tt ON t.id = tt.tag_id
    WHERE tt.thread_id = threads.id
    LIMIT 1
)
WHERE community_id IS NULL 
AND EXISTS (SELECT 1 FROM sqlite_master WHERE type='table' AND name='tags');

-- 6. Add all existing thread authors as members of their communities
INSERT OR IGNORE INTO community_memberships (community_id, user_id, role)
SELECT DISTINCT 
    threads.community_id,
    threads.author_id,
    'member'
FROM threads 
WHERE threads.community_id IS NOT NULL;

-- 7. Add admin as creator/moderator of all auto-created communities
INSERT OR IGNORE INTO community_memberships (community_id, user_id, role)
SELECT id, creator_id, 'creator'
FROM communities
WHERE creator_id = 1; -- Admin user

-- 8. Create default "General" community if no communities exist
INSERT OR IGNORE INTO communities (name, display_name, description, creator_id, visibility, join_approval)
SELECT 'general', 'General', 'General discussion for all topics', 1, 'public', 'open'
WHERE NOT EXISTS (SELECT 1 FROM communities);

-- 9. Add admin as creator of the General community
INSERT OR IGNORE INTO community_memberships (community_id, user_id, role)
SELECT c.id, 1, 'creator'
FROM communities c
WHERE c.name = 'general';

-- 10. Assign orphaned threads (those without community_id) to General community
UPDATE threads 
SET community_id = (SELECT id FROM communities WHERE name = 'general' LIMIT 1)
WHERE community_id IS NULL;

-- 11. Update member counts for all communities
UPDATE communities 
SET member_count = (
    SELECT COUNT(*) 
    FROM community_memberships cm 
    WHERE cm.community_id = communities.id AND cm.status = 'active'
);

-- New migration: Add is_edited column to messages table
ALTER TABLE messages ADD COLUMN is_edited BOOLEAN DEFAULT 0;

COMMIT;

-- Verification queries (run these to check the migration)
-- SELECT 'Communities created:' as check_type, COUNT(*) as count FROM communities;
-- SELECT 'Community memberships:' as check_type, COUNT(*) as count FROM community_memberships;
-- SELECT 'Threads with communities:' as check_type, COUNT(*) as count FROM threads WHERE community_id IS NOT NULL;
-- SELECT 'Orphaned threads:' as check_type, COUNT(*) as count FROM threads WHERE community_id IS NULL; 