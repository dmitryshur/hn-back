package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"time"
)

type Story struct {
	Id          int64    `json:"id"`
	Deleted     *bool    `json:"deleted"`
	Type        string   `json:"type"`
	By          *string  `json:"by"`
	Time        *int     `json:"time"`
	Dead        *bool    `json:"dead"`
	Kids        *[]int64 `json:"kids,omitempty"`
	Descendants *int     `json:"descendants"`
	Score       *int     `json:"score"`
	Title       *string  `json:"title"`
	Url         *string  `json:"url"`
}

type StoryModel struct {
	DB *sql.DB
}

func (s StoryModel) Get(id int64) (*Story, error) {
	query := `SELECT id, deleted, type, by, time, dead, descendants, score, title, url
				FROM stories
				WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var story Story

	err := s.DB.QueryRowContext(ctx, query, id).Scan(
		&story.Id,
		&story.Deleted,
		&story.Type,
		&story.By,
		&story.Time,
		&story.Dead,
		&story.Descendants,
		&story.Score,
		&story.Title,
		&story.Url,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &story, nil
}

func (s StoryModel) GetAll(t string) ([]*Story, error) {
	sort := IfElse(t == "newest", "date", "score")
	query := fmt.Sprintf(`
		SELECT id, deleted, type, by, time, dead, descendants, score, title, url
		FROM stories
		ORDER BY %s DESC`, sort)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var stories []*Story

	for rows.Next() {
		var story Story

		err := rows.Scan(
			&story.Id,
			&story.Deleted,
			&story.Type,
			&story.By,
			&story.Time,
			&story.Dead,
			&story.Descendants,
			&story.Score,
			&story.Title,
			&story.Url,
		)

		if err != nil {
			return nil, err
		}

		stories = append(stories, &story)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stories, nil
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
