package main

type Db interface {
	Save(story *Item, comments []Item) error
}

type Store struct {
	state map[int]int
}

func NewStore() *Store {
	return &Store{state: make(map[int]int)}
}

func (s *Store) Save(story *Item, comments []Item) error {
	return nil
}
