package main

import (
	"github.com/dmitryshur/hackernews/internal/data"
	"net/http"
	"time"
)

const baseUrl = "https://hacker-news.firebaseio.com/v0"

// TODO: need a method to update a story (fetch all the comments)
func main() {
	api := NewApi(http.DefaultClient, baseUrl)
	s := data.NewModel(nil)
	f := NewFetcher(time.Second*600, api, s)

	f.Start()
}
