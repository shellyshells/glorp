package models

import (
	"time"

	"goforum/config"
)

type Message struct {
	ID        int       `json:"id"`
	ThreadID  int       `json:"thread_id"`
	AuthorID  int       `json:"author_id"`
	Author    *User     `json:"author,omitempty"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Likes     int       `json:"likes"`
	Dislikes  int       `json:"dislikes"`
	Score     int       `json:"score"`
	UserVote  int       `json:"user_vote,omitempty"` // 1 for like, -1 for dislike, 0 for no vote
}

type MessageFilters struct {
	ThreadID int    `json:"thread_id"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	SortBy   string `json:"sort_by"` // date, popularity
	UserID   int    `json:"user_id"` // For getting user's vote
}

func CreateMessage(threadID, authorID int, content string) (*Message, error) {
	query := `INSERT INTO messages (thread_id, author_id, content) VALUES (?, ?, ?)`
	result, err := config.DB.Exec(query, threadID, authorID, content)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetMessageByID(int(id), 0)
}

func GetMessageByID(id, userID int) (*Message, error) {
	message := &Message{}
	query := `
		SELECT m.id, m.thread_id, m.author_id, m.content, m.created_at, u.username,
		       COALESCE(SUM(CASE WHEN v.vote_type = 1 THEN 1 ELSE 0 END), 0) as likes,
		       COALESCE(SUM(CASE WHEN v.vote_type = -1 THEN 1 ELSE 0 END), 0) as dislikes,
		       COALESCE(SUM(v.vote_type), 0) as score
		FROM messages m 
		LEFT JOIN users u ON m.author_id = u.id 
		LEFT JOIN votes v ON m.id = v.message_id
		WHERE m.id = ?
		GROUP BY m.id
	`

	var authorUsername string
	err := config.DB.QueryRow(query, id).Scan(
		&message.ID, &message.ThreadID, &message.AuthorID, &message.Content,
		&message.CreatedAt, &authorUsername, &message.Likes, &message.Dislikes,
		&message.Score,
	)

	if err != nil {
		return nil, err
	}

	message.Author = &User{
		ID:       message.AuthorID,
		Username: authorUsername,
	}

	// Get user's vote if userID is provided
	if userID > 0 {
		var vote int
		voteQuery := `SELECT vote_type FROM votes WHERE message_id = ? AND user_id = ?`
		err = config.DB.QueryRow(voteQuery, message.ID, userID).Scan(&vote)
		if err == nil {
			message.UserVote = vote
		}
	}

	return message, nil
}

func GetMessages(filters MessageFilters) ([]Message, int, error) {
	var messages []Message

	// Count total messages
	countQuery := `SELECT COUNT(*) FROM messages WHERE thread_id = ?`
	var total int
	err := config.DB.QueryRow(countQuery, filters.ThreadID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Build main query
	query := `
		SELECT m.id, m.thread_id, m.author_id, m.content, m.created_at, u.username,
		       COALESCE(SUM(CASE WHEN v.vote_type = 1 THEN 1 ELSE 0 END), 0) as likes,
		       COALESCE(SUM(CASE WHEN v.vote_type = -1 THEN 1 ELSE 0 END), 0) as dislikes,
		       COALESCE(SUM(v.vote_type), 0) as score
		FROM messages m 
		LEFT JOIN users u ON m.author_id = u.id 
		LEFT JOIN votes v ON m.id = v.message_id
		WHERE m.thread_id = ?
		GROUP BY m.id
	`

	// Add sorting
	switch filters.SortBy {
	case "popularity":
		query += " ORDER BY score DESC, m.created_at DESC"
	default: // date
		query += " ORDER BY m.created_at ASC"
	}

	// Add pagination
	args := []interface{}{filters.ThreadID}
	if filters.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, filters.Limit)

		if filters.Page > 0 {
			offset := (filters.Page - 1) * filters.Limit
			query += " OFFSET ?"
			args = append(args, offset)
		}
	}

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var message Message
		var authorUsername string

		err := rows.Scan(
			&message.ID, &message.ThreadID, &message.AuthorID, &message.Content,
			&message.CreatedAt, &authorUsername, &message.Likes, &message.Dislikes,
			&message.Score,
		)
		if err != nil {
			continue
		}

		message.Author = &User{
			ID:       message.AuthorID,
			Username: authorUsername,
		}

		// Get user's vote if userID is provided
		if filters.UserID > 0 {
			var vote int
			voteQuery := `SELECT vote_type FROM votes WHERE message_id = ? AND user_id = ?`
			err = config.DB.QueryRow(voteQuery, message.ID, filters.UserID).Scan(&vote)
			if err == nil {
				message.UserVote = vote
			}
		}

		messages = append(messages, message)
	}

	return messages, total, nil
}

func DeleteMessage(messageID int) error {
	// This will also delete related votes due to CASCADE DELETE
	query := `DELETE FROM messages WHERE id = ?`
	_, err := config.DB.Exec(query, messageID)
	return err
}

func CanUserModifyMessage(messageID, userID int, userRole string) bool {
	if userRole == "admin" {
		return true
	}

	var authorID int
	query := `SELECT author_id FROM messages WHERE id = ?`
	err := config.DB.QueryRow(query, messageID).Scan(&authorID)
	if err != nil {
		return false
	}

	return authorID == userID
}

func GetMessageCount() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM messages`
	err := config.DB.QueryRow(query).Scan(&count)
	return count, err
}

func GetMessageCountByUser(userID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM messages WHERE author_id = ?`
	err := config.DB.QueryRow(query, userID).Scan(&count)
	return count, err
}

func GetThreadCountByUser(userID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM threads WHERE author_id = ?`
	err := config.DB.QueryRow(query, userID).Scan(&count)
	return count, err
}

func GetMessagesByUser(userID int, filters MessageFilters) ([]Message, int, error) {
	var messages []Message

	// Count total messages by user
	countQuery := `SELECT COUNT(*) FROM messages WHERE author_id = ?`
	var total int
	err := config.DB.QueryRow(countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Build main query
	query := `
		SELECT m.id, m.thread_id, m.author_id, m.content, m.created_at, u.username,
		       COALESCE(SUM(CASE WHEN v.vote_type = 1 THEN 1 ELSE 0 END), 0) as likes,
		       COALESCE(SUM(CASE WHEN v.vote_type = -1 THEN 1 ELSE 0 END), 0) as dislikes,
		       COALESCE(SUM(v.vote_type), 0) as score
		FROM messages m 
		LEFT JOIN users u ON m.author_id = u.id 
		LEFT JOIN votes v ON m.id = v.message_id
		WHERE m.author_id = ?
		GROUP BY m.id
		ORDER BY m.created_at DESC
	`

	// Add pagination
	args := []interface{}{userID}
	if filters.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, filters.Limit)

		if filters.Page > 0 {
			offset := (filters.Page - 1) * filters.Limit
			query += " OFFSET ?"
			args = append(args, offset)
		}
	}

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var message Message
		var authorUsername string

		err := rows.Scan(
			&message.ID, &message.ThreadID, &message.AuthorID, &message.Content,
			&message.CreatedAt, &authorUsername, &message.Likes, &message.Dislikes,
			&message.Score,
		)
		if err != nil {
			continue
		}

		message.Author = &User{
			ID:       message.AuthorID,
			Username: authorUsername,
		}

		messages = append(messages, message)
	}

	return messages, total, nil
}
