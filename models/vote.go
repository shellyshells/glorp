package models

import (
	"database/sql"
	"time"

	"glorp/config"
)

type Vote struct {
	ID        int       `json:"id"`
	MessageID int       `json:"message_id"`
	UserID    int       `json:"user_id"`
	VoteType  int       `json:"vote_type"` // 1 for like, -1 for dislike
	CreatedAt time.Time `json:"created_at"`
}

func VoteMessage(messageID, userID, voteType int) error {
	// Check if user already voted
	var existingVote int
	checkQuery := `SELECT vote_type FROM votes WHERE message_id = ? AND user_id = ?`
	err := config.DB.QueryRow(checkQuery, messageID, userID).Scan(&existingVote)
	
	if err == sql.ErrNoRows {
		// No existing vote, create new one
		insertQuery := `INSERT INTO votes (message_id, user_id, vote_type) VALUES (?, ?, ?)`
		_, err = config.DB.Exec(insertQuery, messageID, userID, voteType)
		return err
	} else if err != nil {
		return err
	}

	// User already voted
	if existingVote == voteType {
		// Same vote type, remove the vote (toggle off)
		deleteQuery := `DELETE FROM votes WHERE message_id = ? AND user_id = ?`
		_, err = config.DB.Exec(deleteQuery, messageID, userID)
		return err
	} else {
		// Different vote type, update the vote
		updateQuery := `UPDATE votes SET vote_type = ?, created_at = CURRENT_TIMESTAMP 
						WHERE message_id = ? AND user_id = ?`
		_, err = config.DB.Exec(updateQuery, voteType, messageID, userID)
		return err
	}
}

func GetUserVote(messageID, userID int) (int, error) {
	var voteType int
	query := `SELECT vote_type FROM votes WHERE message_id = ? AND user_id = ?`
	err := config.DB.QueryRow(query, messageID, userID).Scan(&voteType)
	if err == sql.ErrNoRows {
		return 0, nil // No vote found
	}
	return voteType, err
}

func GetMessageVoteStats(messageID int) (likes int, dislikes int, score int, err error) {
	query := `
		SELECT 
			COALESCE(SUM(CASE WHEN vote_type = 1 THEN 1 ELSE 0 END), 0) as likes,
			COALESCE(SUM(CASE WHEN vote_type = -1 THEN 1 ELSE 0 END), 0) as dislikes,
			COALESCE(SUM(vote_type), 0) as score
		FROM votes 
		WHERE message_id = ?
	`
	
	err = config.DB.QueryRow(query, messageID).Scan(&likes, &dislikes, &score)
	return
}

func RemoveVotesForMessage(messageID int) error {
	query := `DELETE FROM votes WHERE message_id = ?`
	_, err := config.DB.Exec(query, messageID)
	return err
}

func RemoveVotesForUser(userID int) error {
	query := `DELETE FROM votes WHERE user_id = ?`
	_, err := config.DB.Exec(query, userID)
	return err
}