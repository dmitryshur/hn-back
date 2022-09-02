package data

import "database/sql"

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
	return nil
}
