package main

import (
	"fmt"
	"log"
)
import "github.com/dmitryshur/hackernews/pkg/fetcher"

func main() {
	fmt.Println("hello")
	f := fetcher.NewFetcher()
	stories, err := f.FetchNewestStories()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(stories)
}
