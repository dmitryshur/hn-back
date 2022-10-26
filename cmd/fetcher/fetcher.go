package main

import (
	"github.com/dmitryshur/hackernews/internal/data"
	"github.com/dmitryshur/hackernews/internal/jsonlog"
	"net/http"
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

type config struct {
	fetchInterval time.Duration
	db            struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

type Fetcher struct {
	config config
	logger *jsonlog.Logger
	api    *Api
	store  data.Db
}

type stories []int

func NewFetcher(config config, logger *jsonlog.Logger, api *Api, store data.Db) *Fetcher {
	return &Fetcher{config: config, logger: logger, api: api, store: store}
}

func (f *Fetcher) Start() {
	for {
		bestStories, err := f.FetchBestStories()
		if err != nil {
			f.logger.PrintError(err, nil)
		}
		newestStories, err := f.FetchNewestStories()
		if err != nil {
			f.logger.PrintError(err, nil)
		}

		allStories := stories{}
		allStories = append(bestStories, newestStories...)

		for _, storyId := range allStories {
			story, err := f.FetchItem(storyId)
			if err != nil {
				f.logger.PrintError(err, map[string]string{
					"id": strconv.Itoa(storyId),
				})
			}

			comments, err := f.FetchComments(story)
			if err != nil {
				f.logger.PrintError(err, map[string]string{
					"id": strconv.Itoa(storyId),
				})
			}

			err = f.store.InsertStory(story)
			if err != nil {
				f.logger.PrintError(err, nil)
			}
			f.logger.PrintInfo("inserted story", nil)

			err = f.store.InsertComments(story, *comments)
			if err != nil {
				f.logger.PrintError(err, nil)
			}
			f.logger.PrintInfo("inserted comments", map[string]string{
				"count": strconv.Itoa(len(*comments)),
			})
		}

		if f.config.fetchInterval == 0 {
			break
		}

		time.Sleep(f.config.fetchInterval)
	}
}

func (f *Fetcher) FetchItem(id int) (*data.Item, error) {
	url := strings.Replace(itemUrl, "{{id}}", strconv.Itoa(id), -1)
	url = f.api.baseUrl + url

	response, err := f.api.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var item data.Item
	err = DecodeFromJson(response.Body, &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (f *Fetcher) FetchComments(item *data.Item) (*[]data.Item, error) {
	if item.Kids == nil || len(*item.Kids) == 0 {
		return nil, nil
	}

	commentIds := make(map[int]struct{})
	for _, id := range *item.Kids {
		commentIds[id] = struct{}{}
	}

	var comments []data.Item
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
