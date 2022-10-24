package data

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"time"
)

type Comment struct {
	Id      int     `json:"id"`
	Deleted *bool   `json:"deleted"`
	Type    string  `json:"type"`
	By      *string `json:"by"`
	Time    *int    `json:"time"`
	Dead    *bool   `json:"dead"`
	Kids    *[]int  `json:"kids"`
	Parent  *int    `json:"parent"`
	Text    *string `json:"text"`
	StoryId int     `json:"story_id"`
}

type CommentsModel struct {
	DB *sql.DB
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
			kids = ToPointer([]int{})
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
