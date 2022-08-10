package fetcher

import (
	"net/http"
	"strconv"
	"strings"
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

type Type string

const (
	Story      Type = "story"
	Job        Type = "job"
	Comment    Type = "comment"
	Poll       Type = "poll"
	PollOption Type = "pollopt"
)

type stories []int

type Item struct {
	Id          int     `json:"id"`
	Type        Type    `json:"type"`
	Deleted     *bool   `json:"deleted"`
	By          *string `json:"by"`
	Time        *int    `json:"time"`
	Text        *string `json:"text"`
	Dead        *bool   `json:"dead"`
	Parent      *int    `json:"parent"`
	Kids        *[]int  `json:"kids"`
	Url         *string `json:"url"`
	Score       *int    `json:"score"`
	Title       *string `json:"title"`
	Parts       *[]int  `json:"parts"`
	Descendants *int    `json:"descendants"`
}

func NewFetcher() *Fetcher {
	return &Fetcher{}
}

func (f *Fetcher) FetchItem(id int) (*Item, error) {
	url := strings.Replace(itemUrl, "{{id}}", strconv.Itoa(id), -1)

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var item Item
	err = DecodeFromJson(response.Body, &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
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
