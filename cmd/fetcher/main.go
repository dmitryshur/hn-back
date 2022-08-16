package main

import (
	"net/http"
	"time"
)
import "github.com/dmitryshur/hackernews/pkg/fetcher"

const baseUrl = "https://hacker-news.firebaseio.com/v0"

func main() {
	api := fetcher.NewApi(http.DefaultClient, baseUrl)
	f := fetcher.NewFetcher(time.Second*600, api)

	f.Start()
}
