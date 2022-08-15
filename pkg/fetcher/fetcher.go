package fetcher

import (
	"fmt"
	"github.com/dmitryshur/hackernews/pkg/jsonlog"
	"net/http"
	"os"
	"strconv"
	"strings"
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

func NewFetcher(interval time.Duration, api *Api) *Fetcher {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	return &Fetcher{fetchInterval: interval, logger: logger, api: api}
}

func (f *Fetcher) Start() {
	for {
		bestStories, err := f.FetchBestStories()
		if err != nil {
			f.logger.PrintError(err, nil)
		}
		fmt.Println(bestStories)
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
