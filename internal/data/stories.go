package data

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"time"
)

type Story struct {
	Id          int     `json:"id"`
	Deleted     *bool   `json:"deleted"`
	Type        string  `json:"type"`
	By          *string `json:"by"`
	Time        *int    `json:"time"`
	Dead        *bool   `json:"dead"`
	Kids        *[]int  `json:"kids"`
	Descendants *int    `json:"descendants"`
	Score       *int    `json:"score"`
	Title       *string `json:"title"`
	Url         *string `json:"url"`
}

type StoryModel struct {
	DB *sql.DB
}

func (s StoryModel) Insert(story *Item) error {
	query := `INSERT INTO stories (id, deleted, type, by, time, dead, kids, descendants, score, title, url)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
				ON CONFLICT (id) DO UPDATE 
					SET deleted = $2,
					    type = $3,
					    by = $4,
					    time = $5,
					    dead = $6,
					    kids = $7,
					    descendants = $8,
					    score = $9,
					    title = $10,
					    url = $11`

	args := []interface{}{
		story.Id,
		story.Deleted,
		story.Type,
		story.By,
		story.Time,
		story.Dead,
		pq.Array(*story.Kids),
		story.Descendants,
		story.Score,
		story.Title,
		story.Url,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("insert %w", err)
	}

	return nil
}