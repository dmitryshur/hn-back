package data

import (
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrRecordNotFound = errors.New("record now found")
	ErrEditConflict   = errors.New("edit conflict")
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
	fmt.Println("inserting story")

	return nil
}

func (m Models) InsertComments(story *Item, comments []Item) error {
	fmt.Println("inserting comments")

	return nil
}
