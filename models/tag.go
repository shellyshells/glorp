package models

import (
	"time"

	"goforum/config"
)

type Tag struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func GetAllTags() ([]Tag, error) {
	var tags []Tag
	query := `SELECT id, name, created_at FROM tags ORDER BY name`
	
	rows, err := config.DB.Query(query)
	if err != nil {
		return tags, err
	}
	defer rows.Close()

	for rows.Next() {
		var tag Tag
		err := rows.Scan(&tag.ID, &tag.Name, &tag.CreatedAt)
		if err != nil {
			continue
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func GetTagByID(id int) (*Tag, error) {
	tag := &Tag{}
	query := `SELECT id, name, created_at FROM tags WHERE id = ?`
	
	err := config.DB.QueryRow(query, id).Scan(&tag.ID, &tag.Name, &tag.CreatedAt)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

func GetTagByName(name string) (*Tag, error) {
	tag := &Tag{}
	query := `SELECT id, name, created_at FROM tags WHERE name = ?`
	
	err := config.DB.QueryRow(query, name).Scan(&tag.ID, &tag.Name, &tag.CreatedAt)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

func CreateTag(name string) (*Tag, error) {
	query := `INSERT INTO tags (name) VALUES (?)`
	result, err := config.DB.Exec(query, name)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetTagByID(int(id))
}

func CreateTagIfNotExists(name string) (*Tag, error) {
	// Try to get existing tag first
	tag, err := GetTagByName(name)
	if err == nil {
		return tag, nil
	}

	// Create new tag if it doesn't exist
	return CreateTag(name)
}

func GetTagsWithThreadCount() ([]Tag, error) {
	var tags []Tag
	query := `
		SELECT t.id, t.name, t.created_at, COUNT(tt.thread_id) as thread_count
		FROM tags t 
		LEFT JOIN thread_tags tt ON t.id = tt.tag_id
		LEFT JOIN threads th ON tt.thread_id = th.id AND th.status != 'archived'
		GROUP BY t.id
		ORDER BY thread_count DESC, t.name
	`
	
	rows, err := config.DB.Query(query)
	if err != nil {
		return tags, err
	}
	defer rows.Close()

	for rows.Next() {
		var tag Tag
		var threadCount int
		err := rows.Scan(&tag.ID, &tag.Name, &tag.CreatedAt, &threadCount)
		if err != nil {
			continue
		}
		tags = append(tags, tag)
	}

	return tags, nil
}