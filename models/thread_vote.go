package models

import (
	"database/sql"
	"time"

	"glorp/config"
)

type ThreadVote struct {
	ID        int       `json:"id"`
	ThreadID  int       `json:"thread_id"`
	UserID    int       `json:"user_id"`
	VoteType  int       `json:"vote_type"` // 1 for upvote, -1 for downvote
	CreatedAt time.Time `json:"created_at"`
}

func VoteThread(threadID, userID, voteType int) error {
	// Check if user already voted
	var existingVote int
	checkQuery := `SELECT vote_type FROM thread_votes WHERE thread_id = ? AND user_id = ?`
	err := config.DB.QueryRow(checkQuery, threadID, userID).Scan(&existingVote)

	if err == sql.ErrNoRows {
		// No existing vote, create new one
		insertQuery := `INSERT INTO thread_votes (thread_id, user_id, vote_type) VALUES (?, ?, ?)`
		_, err = config.DB.Exec(insertQuery, threadID, userID, voteType)
		return err
	} else if err != nil {
		return err
	}

	// User already voted
	if existingVote == voteType {
		// Same vote type, remove the vote (toggle off)
		deleteQuery := `DELETE FROM thread_votes WHERE thread_id = ? AND user_id = ?`
		_, err = config.DB.Exec(deleteQuery, threadID, userID)
		return err
	} else {
		// Different vote type, update the vote
		updateQuery := `UPDATE thread_votes SET vote_type = ?, created_at = CURRENT_TIMESTAMP 
						WHERE thread_id = ? AND user_id = ?`
		_, err = config.DB.Exec(updateQuery, voteType, threadID, userID)
		return err
	}
}

func GetUserThreadVote(threadID, userID int) (int, error) {
	var voteType int
	query := `SELECT vote_type FROM thread_votes WHERE thread_id = ? AND user_id = ?`
	err := config.DB.QueryRow(query, threadID, userID).Scan(&voteType)
	if err == sql.ErrNoRows {
		return 0, nil // No vote found
	}
	return voteType, err
}

func GetThreadVoteStats(threadID int) (upvotes int, downvotes int, score int, err error) {
	query := `
		SELECT 
			COALESCE(SUM(CASE WHEN vote_type = 1 THEN 1 ELSE 0 END), 0) as upvotes,
			COALESCE(SUM(CASE WHEN vote_type = -1 THEN 1 ELSE 0 END), 0) as downvotes,
			COALESCE(SUM(vote_type), 0) as score
		FROM thread_votes 
		WHERE thread_id = ?
	`

	err = config.DB.QueryRow(query, threadID).Scan(&upvotes, &downvotes, &score)
	return
}
