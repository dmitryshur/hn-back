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

// TODO: add get for all comments
type Db interface {
	GetStory(id int64) (*Story, error)
	GetStories(t string) ([]*Story, error)
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

func (m Models) GetStory(id int64) (*Story, error) {
	story, err := m.Stories.Get(id)
	if err != nil {
		return nil, fmt.Errorf("getStory %w", err)
	}

	return story, nil
}

func (m Models) GetStories(t string) ([]*Story, error) {
	stories, err := m.Stories.GetAll(t)
	if err != nil {
		return nil, fmt.Errorf("getStories %w", err)
	}

	return stories, nil
}

func (m Models) InsertStory(story *Item) error {
	err := m.Stories.Insert(story)
	if err != nil {
		return fmt.Errorf("insertStory %w", err)
	}

	return nil
}

func (m Models) InsertComments(story *Item, comments []Item) error {
	err := m.Comments.Insert(story, comments)
	if err != nil {
		return fmt.Errorf("insertComments %w", err)
	}

	return nil
}
