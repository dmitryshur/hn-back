package main

import (
	"fmt"
	"log"
)
import "github.com/dmitryshur/hackernews/pkg/fetcher"

func main() {
	fmt.Println("hello")
	f := fetcher.NewFetcher()
	item, err := f.FetchItem(8863)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(item)
}
