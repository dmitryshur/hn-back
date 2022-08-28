package main

import (
	"fmt"
)

type Db interface {
	Save(items *[]Item) error
}

type Store struct {
	state map[int]int
}

func NewStore() *Store {
	return &Store{state: make(map[int]int)}
}

func (s *Store) Save(items *[]Item) {
	fmt.Println(len(*items))
}
