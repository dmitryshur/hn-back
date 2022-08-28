package main

import (
	"github.com/dmitryshur/hackernews/internal/jsonlog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	itemUrl            = "/item/{{id}}.json"
	bestStoriesUrl     = "/beststories.json"
	newestStoriesUrl   = "/newstories.json"
	bestStoriesCount   = 200
	newestStoriesCount = 500
)

type Api struct {
	client  *http.Client
	baseUrl string
}

func NewApi(client *http.Client, baseUrl string) *Api {
	return &Api{client: client, baseUrl: baseUrl}
}

type Fetcher struct {
	fetchInterval time.Duration
	logger        *jsonlog.Logger
	stories       map[int]struct{}
	api           *Api
	store         *Store
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
	Poll        *int    `json:"poll"`
	Kids        *[]int  `json:"kids"`
	Url         *string `json:"url"`
	Score       *int    `json:"score"`
	Title       *string `json:"title"`
	Parts       *[]int  `json:"parts"`
	Descendants *int    `json:"descendants"`
}

func NewFetcher(interval time.Duration, api *Api, store *Store) *Fetcher {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	return &Fetcher{fetchInterval: interval, logger: logger, api: api, store: store}
}

// TODO: fetch newest stories. handle comments
func (f *Fetcher) Start() {
	for {
		bestStories, err := f.FetchBestStories()
		if err != nil {
			f.logger.PrintError(err, nil)
		}

		for _, storyId := range bestStories {
			item, err := f.FetchItem(storyId)
			if err != nil {
				f.logger.PrintError(err, map[string]string{
					"id": strconv.Itoa(storyId),
				})
			}

			comments, err := f.FetchComments(item)
			if err != nil {
				f.logger.PrintError(err, map[string]string{
					"id": strconv.Itoa(storyId),
				})
			}
			f.store.Save(comments)
		}

		time.Sleep(f.fetchInterval)
	}
}

func (f *Fetcher) FetchItem(id int) (*Item, error) {
	url := strings.Replace(itemUrl, "{{id}}", strconv.Itoa(id), -1)
	url = f.api.baseUrl + url

	response, err := f.api.client.Get(url)
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

func (f *Fetcher) FetchComments(item *Item) (*[]Item, error) {
	if item.Kids == nil || len(*item.Kids) == 0 {
		return nil, nil
	}

	commentIds := make(map[int]struct{})
	for _, id := range *item.Kids {
		commentIds[id] = struct{}{}
	}

	var comments []Item
	var mu sync.Mutex
	var wg sync.WaitGroup
	for len(commentIds) > 0 {
		for id := range commentIds {
			wg.Add(1)

			go func(id int) {
				defer wg.Done()

				comment, err := f.FetchItem(id)
				if err != nil {
					f.logger.PrintError(err, map[string]string{
						"id": strconv.Itoa(id),
					})
				}

				mu.Lock()
				comments = append(comments, *comment)
				delete(commentIds, comment.Id)
				if comment.Kids != nil && len(*comment.Kids) > 0 {
					for _, kidId := range *comment.Kids {
						commentIds[kidId] = struct{}{}
					}
				}
				mu.Unlock()
			}(id)
		}
		wg.Wait()
	}

	return &comments, nil
}

func (f *Fetcher) FetchBestStories() (stories, error) {
	url := f.api.baseUrl + bestStoriesUrl
	response, err := f.api.client.Get(url)

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
	url := f.api.baseUrl + newestStoriesUrl
	response, err := f.api.client.Get(url)

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
