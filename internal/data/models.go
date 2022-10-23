package data

import (
	"database/sql"
	"fmt"
)

type Db interface {
	InsertStory(story *Item) error

	InsertComments(story *Item, comments []Item) error
}

type Models struct {
	Stories  StoryModel
	Comments CommentsModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		Stories:  StoryModel{DB: db},
		Comments: CommentsModel{DB: db},
	}
}

func (m Models) InsertStory(story *Item) error {
	err := m.Stories.Insert(story)
	if err != nil {
		return fmt.Errorf("insertStory %w", err)
	}

	return nil
}

func (m Models) InsertComments(story *Item, comments []Item) error {
	fmt.Println("inserting comments")

	return nil
}
