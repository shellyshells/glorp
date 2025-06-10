package models

import (
	"database/sql"
	"strings"
	"time"

	"glorp/config"
)

type Thread struct {
	ID           int        `json:"id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	AuthorID     int        `json:"author_id"`
	Author       *User      `json:"author,omitempty"`
	CommunityID  *int       `json:"community_id,omitempty"`
	Community    *Community `json:"community,omitempty"`
	Status       string     `json:"status"`              // open, closed, archived
	PostType     string     `json:"post_type"`           // text, link, image
	ImageURL     string     `json:"image_url,omitempty"` // For image posts
	LinkURL      string     `json:"link_url,omitempty"`  // For link posts
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Tags         []Tag      `json:"tags,omitempty"`
	MessageCount int        `json:"message_count"`
	Upvotes      int        `json:"upvotes"`
	Downvotes    int        `json:"downvotes"`
	Score        int        `json:"score"`
	UserVote     int        `json:"user_vote,omitempty"` // 1 for upvote, -1 for downvote, 0 for no vote
}

type ThreadFilters struct {
	TagID       int    `json:"tag_id"`
	TagName     string `json:"tag_name"`
	Search      string `json:"search"`
	Status      string `json:"status"`
	AuthorID    int    `json:"author_id"`
	UserID      int    `json:"user_id"`
	CommunityID int    `json:"community_id"`
	Page        int    `json:"page"`
	Limit       int    `json:"limit"`
	SortBy      string `json:"sort_by"` // date, popularity, hot
}

func CreateThread(title, description string, authorID int, tagIDs []int) (*Thread, error) {
	return CreateThreadWithMedia(title, description, authorID, tagIDs, "text", "", "")
}

func CreateThreadWithMedia(title, description string, authorID int, tagIDs []int, postType, imageURL, linkURL string) (*Thread, error) {
	tx, err := config.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `INSERT INTO threads (title, description, author_id, post_type, image_url, link_url) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := tx.Exec(query, title, description, authorID, postType, imageURL, linkURL)
	if err != nil {
		return nil, err
	}

	threadID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Add tags
	for _, tagID := range tagIDs {
		_, err = tx.Exec(`INSERT INTO thread_tags (thread_id, tag_id) VALUES (?, ?)`,
			threadID, tagID)
		if err != nil {
			return nil, err
		}
	}

	// Auto-upvote the thread by the author
	_, err = tx.Exec(`INSERT INTO thread_votes (thread_id, user_id, vote_type) VALUES (?, ?, 1)`,
		threadID, authorID)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return GetThreadByIDWithUser(int(threadID), authorID)
}

func GetThreadByID(id int) (*Thread, error) {
	return GetThreadByIDWithUser(id, 0)
}

func GetThreadByIDWithUser(id, userID int) (*Thread, error) {
	thread := &Thread{}
	query := `SELECT t.id, t.title, t.description, t.author_id, t.status, 
			         t.post_type, COALESCE(t.image_url, '') as image_url, COALESCE(t.link_url, '') as link_url,
			         t.created_at, t.updated_at, u.username,
					 COUNT(DISTINCT m.id) as message_count,
					 COALESCE(SUM(CASE WHEN tv.vote_type = 1 THEN 1 ELSE 0 END), 0) as upvotes,
					 COALESCE(SUM(CASE WHEN tv.vote_type = -1 THEN 1 ELSE 0 END), 0) as downvotes,
					 COALESCE(SUM(tv.vote_type), 0) as score
			  FROM threads t 
			  LEFT JOIN users u ON t.author_id = u.id 
			  LEFT JOIN messages m ON t.id = m.thread_id
			  LEFT JOIN thread_votes tv ON t.id = tv.thread_id
			  WHERE t.id = ?
			  GROUP BY t.id`

	var authorUsername string
	err := config.DB.QueryRow(query, id).Scan(
		&thread.ID, &thread.Title, &thread.Description, &thread.AuthorID,
		&thread.Status, &thread.PostType, &thread.ImageURL, &thread.LinkURL,
		&thread.CreatedAt, &thread.UpdatedAt, &authorUsername,
		&thread.MessageCount, &thread.Upvotes, &thread.Downvotes, &thread.Score,
	)

	if err != nil {
		return nil, err
	}

	thread.Author = &User{
		ID:       thread.AuthorID,
		Username: authorUsername,
	}

	// Get user's vote if userID is provided
	if userID > 0 {
		var vote int
		voteQuery := `SELECT vote_type FROM thread_votes WHERE thread_id = ? AND user_id = ?`
		err = config.DB.QueryRow(voteQuery, thread.ID, userID).Scan(&vote)
		if err == nil {
			thread.UserVote = vote
		}
	}

	// Get tags
	thread.Tags, _ = GetThreadTags(thread.ID)

	return thread, nil
}

func GetThreads(filters ThreadFilters) ([]Thread, int, error) {
	var threads []Thread
	var whereConditions []string
	var args []interface{}
	argIndex := 0

	baseQuery := `
		FROM threads t 
		LEFT JOIN users u ON t.author_id = u.id 
		LEFT JOIN messages m ON t.id = m.thread_id
		LEFT JOIN thread_tags tt ON t.id = tt.thread_id
		LEFT JOIN tags tag ON tt.tag_id = tag.id
		LEFT JOIN thread_votes tv ON t.id = tv.thread_id
	`

	// Build WHERE conditions
	if filters.TagID > 0 {
		whereConditions = append(whereConditions, "tt.tag_id = ?")
		args = append(args, filters.TagID)
		argIndex++
	}

	if filters.TagName != "" {
		whereConditions = append(whereConditions, "tag.name = ?")
		args = append(args, filters.TagName)
		argIndex++
	}

	if filters.CommunityID > 0 {
		whereConditions = append(whereConditions, "t.community_id = ?")
		args = append(args, filters.CommunityID)
		argIndex++
	}

	if filters.Search != "" {
		whereConditions = append(whereConditions, "(t.title LIKE ? OR t.description LIKE ?)")
		searchTerm := "%" + filters.Search + "%"
		args = append(args, searchTerm, searchTerm)
		argIndex += 2
	}

	if filters.Status != "" {
		whereConditions = append(whereConditions, "t.status = ?")
		args = append(args, filters.Status)
		argIndex++
	} else {
		// Default: don't show archived threads
		whereConditions = append(whereConditions, "t.status != 'archived'")
	}

	if filters.AuthorID > 0 {
		whereConditions = append(whereConditions, "t.author_id = ?")
		args = append(args, filters.AuthorID)
		argIndex++
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Count total
	countQuery := "SELECT COUNT(DISTINCT t.id) " + baseQuery + whereClause
	var total int
	err := config.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Main query with voting data
	selectQuery := `
		SELECT DISTINCT t.id, t.title, t.description, t.author_id, t.status, 
		       t.post_type, COALESCE(t.image_url, '') as image_url, COALESCE(t.link_url, '') as link_url,
		       t.created_at, t.updated_at, u.username,
			   COUNT(DISTINCT m.id) as message_count,
			   COALESCE(SUM(CASE WHEN tv.vote_type = 1 THEN 1 ELSE 0 END), 0) as upvotes,
			   COALESCE(SUM(CASE WHEN tv.vote_type = -1 THEN 1 ELSE 0 END), 0) as downvotes,
			   COALESCE(SUM(tv.vote_type), 0) as score
	` + baseQuery + whereClause + ` 
		GROUP BY t.id
	`

	// Add sorting
	switch filters.SortBy {
	case "hot":
		// Hot algorithm: score / (age in hours + 2)^1.8
		selectQuery += " ORDER BY (COALESCE(SUM(tv.vote_type), 0) + 1) / POW((julianday('now') - julianday(t.created_at)) * 24 + 2, 1.8) DESC"
	case "top":
		selectQuery += " ORDER BY COALESCE(SUM(tv.vote_type), 0) DESC, t.created_at DESC"
	case "new":
		selectQuery += " ORDER BY t.created_at DESC"
	default: // date or popularity
		selectQuery += " ORDER BY t.created_at DESC"
	}

	// Add pagination
	if filters.Limit > 0 {
		selectQuery += " LIMIT ?"
		args = append(args, filters.Limit)

		if filters.Page > 0 {
			offset := (filters.Page - 1) * filters.Limit
			selectQuery += " OFFSET ?"
			args = append(args, offset)
		}
	}

	rows, err := config.DB.Query(selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var thread Thread
		var authorUsername string

		err := rows.Scan(
			&thread.ID, &thread.Title, &thread.Description, &thread.AuthorID,
			&thread.Status, &thread.PostType, &thread.ImageURL, &thread.LinkURL,
			&thread.CreatedAt, &thread.UpdatedAt, &authorUsername,
			&thread.MessageCount, &thread.Upvotes, &thread.Downvotes, &thread.Score,
		)
		if err != nil {
			continue
		}

		thread.Author = &User{
			ID:       thread.AuthorID,
			Username: authorUsername,
		}

		// Get user's vote if userID is provided
		if filters.UserID > 0 {
			var vote int
			voteQuery := `SELECT vote_type FROM thread_votes WHERE thread_id = ? AND user_id = ?`
			err = config.DB.QueryRow(voteQuery, thread.ID, filters.UserID).Scan(&vote)
			if err == nil {
				thread.UserVote = vote
			}
		}

		// Get tags for each thread
		thread.Tags, _ = GetThreadTags(thread.ID)

		threads = append(threads, thread)
	}

	return threads, total, nil
}

func UpdateThread(threadID int, title, description string, tagIDs []int) error {
	return UpdateThreadWithMedia(threadID, title, description, "", "", tagIDs)
}

func UpdateThreadWithMedia(threadID int, title, description, imageURL, linkURL string, tagIDs []int) error {
	tx, err := config.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update thread
	query := `UPDATE threads SET title = ?, description = ?, image_url = ?, link_url = ?, updated_at = CURRENT_TIMESTAMP 
			  WHERE id = ?`
	_, err = tx.Exec(query, title, description, imageURL, linkURL, threadID)
	if err != nil {
		return err
	}

	// Remove existing tags
	_, err = tx.Exec(`DELETE FROM thread_tags WHERE thread_id = ?`, threadID)
	if err != nil {
		return err
	}

	// Add new tags
	for _, tagID := range tagIDs {
		_, err = tx.Exec(`INSERT INTO thread_tags (thread_id, tag_id) VALUES (?, ?)`,
			threadID, tagID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func DeleteThread(threadID int) error {
	// Due to CASCADE DELETE, this will also delete related messages, votes, and thread_tags
	query := `DELETE FROM threads WHERE id = ?`
	_, err := config.DB.Exec(query, threadID)
	return err
}

func UpdateThreadStatus(threadID int, status string) error {
	query := `UPDATE threads SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := config.DB.Exec(query, status, threadID)
	return err
}

func GetThreadTags(threadID int) ([]Tag, error) {
	var tags []Tag
	query := `SELECT t.id, t.name FROM tags t 
			  JOIN thread_tags tt ON t.id = tt.tag_id 
			  WHERE tt.thread_id = ?`

	rows, err := config.DB.Query(query, threadID)
	if err != nil {
		return tags, err
	}
	defer rows.Close()

	for rows.Next() {
		var tag Tag
		err := rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			continue
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func CanUserModifyThread(threadID, userID int, userRole string) bool {
	if userRole == "admin" {
		return true
	}

	var authorID int
	query := `SELECT author_id FROM threads WHERE id = ?`
	err := config.DB.QueryRow(query, threadID).Scan(&authorID)
	if err != nil {
		return false
	}

	return authorID == userID
}

// Helper to get threads by community
func GetThreadsByCommunity(filters ThreadFilters) ([]Thread, int, error) {
	filters.CommunityID = filters.CommunityID // ensure field is set
	return GetThreads(filters)
}

// CreateThreadInCommunity creates a thread in a specific community
func CreateThreadInCommunity(title, description string, authorID, communityID int, postType, imageURL, linkURL string) (*Thread, error) {
	tx, err := config.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `INSERT INTO threads (title, description, author_id, community_id, post_type, image_url, link_url) VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := tx.Exec(query, title, description, authorID, communityID, postType, imageURL, linkURL)
	if err != nil {
		return nil, err
	}

	threadID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Auto-upvote the thread by the author
	_, err = tx.Exec(`INSERT INTO thread_votes (thread_id, user_id, vote_type) VALUES (?, ?, 1)`,
		threadID, authorID)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return GetThreadByIDWithUserAndCommunity(int(threadID), authorID)
}

// Helper function to check if user can post in community
func CanUserPostInCommunity(communityID, userID int) (bool, error) {
	// Get community info
	var visibility, joinApproval string
	err := config.DB.QueryRow("SELECT visibility, join_approval FROM communities WHERE id = ?", communityID).Scan(&visibility, &joinApproval)
	if err != nil {
		return false, err
	}

	// If public community with open posting, anyone can post
	if visibility == "public" && joinApproval == "open" {
		return true, nil
	}

	// Check if user is a member
	var role string
	err = config.DB.QueryRow(`
		SELECT role FROM community_memberships 
		WHERE community_id = ? AND user_id = ? AND status = 'active'`,
		communityID, userID).Scan(&role)

	if err != nil {
		// User is not a member
		if visibility == "private" {
			return false, nil // Private communities require membership
		}
		if visibility == "restricted" {
			return false, nil // Restricted communities require membership to post
		}
		return true, nil // Public communities allow non-members to view/post
	}

	// User is a member, so they can post
	return true, nil
}

// Update GetThreadByIDWithUser to include community info
func GetThreadByIDWithUserAndCommunity(id, userID int) (*Thread, error) {
	thread := &Thread{}
	query := `SELECT t.id, t.title, t.description, t.author_id, t.community_id, t.status, 
			         t.post_type, COALESCE(t.image_url, '') as image_url, COALESCE(t.link_url, '') as link_url,
			         t.created_at, t.updated_at, u.username, c.name as community_name, c.display_name as community_display_name,
					 COUNT(DISTINCT m.id) as message_count,
					 COALESCE(SUM(CASE WHEN tv.vote_type = 1 THEN 1 ELSE 0 END), 0) as upvotes,
					 COALESCE(SUM(CASE WHEN tv.vote_type = -1 THEN 1 ELSE 0 END), 0) as downvotes,
					 COALESCE(SUM(tv.vote_type), 0) as score
			  FROM threads t 
			  LEFT JOIN users u ON t.author_id = u.id 
			  LEFT JOIN messages m ON t.id = m.thread_id
			  LEFT JOIN thread_votes tv ON t.id = tv.thread_id
			  LEFT JOIN communities c ON t.community_id = c.id
			  WHERE t.id = ?
			  GROUP BY t.id`

	var authorUsername, communityName, communityDisplayName sql.NullString
	var communityID sql.NullInt64
	err := config.DB.QueryRow(query, id).Scan(
		&thread.ID, &thread.Title, &thread.Description, &thread.AuthorID, &communityID,
		&thread.Status, &thread.PostType, &thread.ImageURL, &thread.LinkURL,
		&thread.CreatedAt, &thread.UpdatedAt, &authorUsername, &communityName, &communityDisplayName,
		&thread.MessageCount, &thread.Upvotes, &thread.Downvotes, &thread.Score,
	)

	if err != nil {
		return nil, err
	}

	thread.Author = &User{
		ID:       thread.AuthorID,
		Username: authorUsername.String,
	}

	// Set community info if thread belongs to a community
	if communityID.Valid {
		thread.CommunityID = new(int)
		*thread.CommunityID = int(communityID.Int64)
		thread.Community = &Community{
			ID:          int(communityID.Int64),
			Name:        communityName.String,
			DisplayName: communityDisplayName.String,
		}
	}

	// Get user's vote if userID is provided
	if userID > 0 {
		var vote int
		voteQuery := `SELECT vote_type FROM thread_votes WHERE thread_id = ? AND user_id = ?`
		err = config.DB.QueryRow(voteQuery, thread.ID, userID).Scan(&vote)
		if err == nil {
			thread.UserVote = vote
		}
	}

	// Get tags for backward compatibility
	thread.Tags, _ = GetThreadTags(thread.ID)

	return thread, nil
}
