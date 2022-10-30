package data

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"time"
)

type Comment struct {
	Id      int64    `json:"id"`
	Deleted *bool    `json:"deleted"`
	Type    string   `json:"type"`
	By      *string  `json:"by"`
	Time    *int     `json:"time"`
	Dead    *bool    `json:"dead"`
	Kids    *[]int64 `json:"kids"`
	Parent  *int     `json:"parent"`
	Text    *string  `json:"text"`
	StoryId int      `json:"story_id"`
}

type CommentsModel struct {
	DB *sql.DB
}

func (c CommentsModel) GetAll(storyId int64) ([]*Comment, error) {
	query := fmt.Sprintf(`
		SELECT comments.id, comments.deleted, comments.type, comments.by, comments.time, comments.dead, comments.kids, parent, story_id, text
		FROM comments
		INNER JOIN stories
		ON stories.id = comments.story_id
		WHERE story_id = '%d'
		ORDER BY parent`, storyId)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := c.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var comments []*Comment

	for rows.Next() {
		var comment Comment

		comment.Kids = ToPointer([]int64{})
		err := rows.Scan(
			&comment.Id,
			&comment.Deleted,
			&comment.Type,
			&comment.By,
			&comment.Time,
			&comment.Dead,
			pq.Array(comment.Kids),
			&comment.Parent,
			&comment.StoryId,
			&comment.Text,
		)

		if err != nil {
			return nil, err
		}

		comments = append(comments, &comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (c CommentsModel) Insert(story *Item, comments []Item) error {
	// FIXME: bulk insertion is needed
	query := `INSERT INTO comments (id, deleted, type, by, time, dead, kids, parent, story_id, text)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
				ON CONFLICT (id) DO UPDATE
					SET deleted = $2,
						type = $3,
						by = $4,
						time = $5,
						dead = $6,
						kids = $7,
						parent = $8,
						story_id = $9,
						text = $10`

	ctx, cancel := context.WithCancel(context.Background())
	timer := time.AfterFunc(5*time.Second, cancel)
	defer timer.Stop()

	for _, comment := range comments {
		timer.Reset(5 * time.Second)

		kids := comment.Kids
		if kids == nil {
			kids = ToPointer([]int64{})
		}

		args := []interface{}{
			comment.Id,
			comment.Deleted,
			comment.Type,
			comment.By,
			comment.Time,
			comment.Dead,
			pq.Array(*kids),
			comment.Parent,
			story.Id,
			comment.Text,
		}

		_, err := c.DB.ExecContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("insert %w", err)
		}
	}

	return nil
}
