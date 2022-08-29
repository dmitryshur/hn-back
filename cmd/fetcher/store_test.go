package main

import (
	"testing"
)

type MockStore struct {
	state map[int][]Item
}

func NewMockStore() *MockStore {
	return &MockStore{
		state: make(map[int][]Item),
	}
}

func (s *MockStore) Save(story *Item, comments []Item) error {
	s.state[story.Id] = comments
	return nil
}

func TestStore(t *testing.T) {
	t.Run("save story with a single comment", func(t *testing.T) {
		store := NewMockStore()
		story := Item{
			Id:   32411232,
			Type: "story",
			Kids: &[]int{
				32411282,
			},
		}

		comments := []Item{
			{
				Id:   32411282,
				Type: "comment",
			},
		}
		expected := map[int][]Item{
			32411232: {
				{
					Id:   32411282,
					Type: "comment",
				},
			},
		}

		store.Save(&story, comments)
		for storyId, comments := range store.state {
			if _, ok := expected[storyId]; !ok {
				t.Errorf("missing comments for story %d in state", storyId)
			}

			for _, comment := range expected[storyId] {
				if !Includes(comments, func(c Item) bool {
					return comment.Id == c.Id
				}) {
					t.Errorf("expected %v to be in state", comment)
				}
			}
		}
	})
}
