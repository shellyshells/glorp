package models

import (
	"strings"
	"time"

	"goforum/config"
)

type Thread struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	AuthorID    int       `json:"author_id"`
	Author      *User     `json:"author,omitempty"`
	Status      string    `json:"status"` // open, closed, archived
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Tags        []Tag     `json:"tags,omitempty"`
	MessageCount int      `json:"message_count"`
}

type ThreadFilters struct {
	TagID    int    `json:"tag_id"`
	TagName  string `json:"tag_name"`
	Search   string `json:"search"`
	Status   string `json:"status"`
	AuthorID int    `json:"author_id"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	SortBy   string `json:"sort_by"` // date, popularity
}

func CreateThread(title, description string, authorID int, tagIDs []int) (*Thread, error) {
	tx, err := config.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `INSERT INTO threads (title, description, author_id) VALUES (?, ?, ?)`
	result, err := tx.Exec(query, title, description, authorID)
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

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return GetThreadByID(int(threadID))
}

func GetThreadByID(id int) (*Thread, error) {
	thread := &Thread{}
	query := `SELECT t.id, t.title, t.description, t.author_id, t.status, 
			         t.created_at, t.updated_at, u.username,
					 COUNT(m.id) as message_count
			  FROM threads t 
			  LEFT JOIN users u ON t.author_id = u.id 
			  LEFT JOIN messages m ON t.id = m.thread_id
			  WHERE t.id = ?
			  GROUP BY t.id`
	
	var authorUsername string
	err := config.DB.QueryRow(query, id).Scan(
		&thread.ID, &thread.Title, &thread.Description, &thread.AuthorID,
		&thread.Status, &thread.CreatedAt, &thread.UpdatedAt, &authorUsername,
		&thread.MessageCount,
	)
	
	if err != nil {
		return nil, err
	}

	thread.Author = &User{
		ID:       thread.AuthorID,
		Username: authorUsername,
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

	// Main query
	selectQuery := `
		SELECT DISTINCT t.id, t.title, t.description, t.author_id, t.status, 
		       t.created_at, t.updated_at, u.username,
			   COUNT(DISTINCT m.id) as message_count
	` + baseQuery + whereClause + ` 
		GROUP BY t.id
		ORDER BY t.created_at DESC
	`

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
			&thread.Status, &thread.CreatedAt, &thread.UpdatedAt, &authorUsername,
			&thread.MessageCount,
		)
		if err != nil {
			continue
		}

		thread.Author = &User{
			ID:       thread.AuthorID,
			Username: authorUsername,
		}

		// Get tags for each thread
		thread.Tags, _ = GetThreadTags(thread.ID)

		threads = append(threads, thread)
	}

	return threads, total, nil
}

func UpdateThread(threadID int, title, description string, tagIDs []int) error {
	tx, err := config.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update thread
	query := `UPDATE threads SET title = ?, description = ?, updated_at = CURRENT_TIMESTAMP 
			  WHERE id = ?`
	_, err = tx.Exec(query, title, description, threadID)
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