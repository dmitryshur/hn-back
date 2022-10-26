package data

import (
	"database/sql"
	"fmt"
)

type Db interface {
	InsertStory(story *Item) error
	GetStories(t string) ([]*Story, error)

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

func (m Models) GetStories(t string) ([]*Story, error) {
	if t == "newest" {
		stories, err := m.Stories.GetNewest()
		if err != nil {
			return nil, fmt.Errorf("newest getStories %w", err)
		}

		return stories, nil
	} else {
		err := m.Stories.GetBest()
		if err != nil {
			return nil, fmt.Errorf("best getStories %w", err)
		}
	}

	return nil, nil
}

func (m Models) InsertComments(story *Item, comments []Item) error {
	err := m.Comments.Insert(story, comments)
	if err != nil {
		return fmt.Errorf("insertComments %w", err)
	}

	return nil
}
