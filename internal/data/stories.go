package data

import "database/sql"

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

func (s StoryModel) Insert(story *Item, comments []Item) error {
	return nil
}
