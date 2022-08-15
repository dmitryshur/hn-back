package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)
import "github.com/dmitryshur/hackernews/pkg/fetcher"

const baseUrl = "https://hacker-news.firebaseio.com/v0"

// TODO: create fn to work with the response of best/newest stories. it should go through all the ids, get each story
// TODO: for each story, fetch the comments as well
func main() {
	api := fetcher.NewApi(http.DefaultClient, baseUrl)
	f := fetcher.NewFetcher(time.Second*5, api)
	stories, err := f.FetchItem(8863)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(stories)
}
