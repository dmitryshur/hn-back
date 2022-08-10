package fetcher

import (
	"net/http"
)

const (
	itemUrl            = "https://hacker-news.firebaseio.com/v0/item/{{id}}.json"
	bestStoriesUrl     = "https://hacker-news.firebaseio.com/v0/beststories.json"
	newestStoriesUrl   = "https://hacker-news.firebaseio.com/v0/newstories.json"
	bestStoriesCount   = 200
	newestStoriesCount = 500
)

type Fetcher struct {
}

type stories []int

func NewFetcher() *Fetcher {
	return &Fetcher{}
}

func (f *Fetcher) FetchItem() {

}

func (f *Fetcher) FetchBestStories() (stories, error) {
	response, err := http.Get(bestStoriesUrl)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var input = make([]int, bestStoriesCount)
	err = DecodeFromJson(response.Body, &input)
	if err != nil {
		return nil, err
	}

	return input, nil
}

func (f *Fetcher) FetchNewestStories() (stories, error) {
	response, err := http.Get(newestStoriesUrl)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var input = make([]int, newestStoriesCount)
	err = DecodeFromJson(response.Body, &input)
	if err != nil {
		return nil, err
	}

	return input, nil
}
