package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record now found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Db interface {
	InsertStory(story *Item) error
	GetStory() error

	InsertComments(story *Item, comments []Item) error
	GetComments() error
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
	panic("implement me")
}

func (m Models) GetStory() error {
	panic("implement me")
}

func (m Models) InsertComments(story *Item, comments []Item) error {
	panic("implement me")
}

func (m Models) GetComments() error {
	panic("implement me")
}
