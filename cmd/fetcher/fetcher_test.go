package main

import (
	"github.com/dmitryshur/hackernews/internal/data"
	"github.com/dmitryshur/hackernews/internal/jsonlog"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
)

type ModelsMock struct {
	state map[int][]int
}

func NewModelsMock() *ModelsMock {
	return &ModelsMock{state: map[int][]int{}}
}

func (m ModelsMock) InsertStory(story *data.Item) error {
	if _, ok := m.state[story.Id]; !ok {
		m.state[story.Id] = []int{}
	}

	return nil
}

// TODO: add tests
func (m ModelsMock) GetStories(t string) ([]*data.Story, error) {
	return nil, nil
}

func (m ModelsMock) InsertComments(story *data.Item, comments []data.Item) error {
	if _, ok := m.state[story.Id]; !ok {
		m.state[story.Id] = []int{}
	}

	commentsIds := make([]int, len(comments))
	for _, comment := range comments {
		commentsIds = append(commentsIds, comment.Id)
	}

	m.state[story.Id] = append(m.state[story.Id], commentsIds...)

	return nil
}

func TestFetcher(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mocker := NewMocker()

		switch r.URL.String() {
		case "/beststories.json", "/newstories.json":
			w.Write(mocker.Stories())
		default:
			url := r.URL.String()
			id, err := strconv.Atoi(url[6 : len(url)-5])
			if err != nil {
				t.Errorf("Can't convert id to int. %s", err)
			}

			w.Write(mocker.Item(id))
		}
	}))
	defer server.Close()

	t.Run("fetch best stories", func(t *testing.T) {
		api := NewApi(server.Client(), server.URL)
		store := NewModelsMock()
		config := config{fetchInterval: time.Second * 5}
		logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

		fetcher := NewFetcher(config, logger, api, store)
		expected := []int{
			32411232,
			32627286,
			32626745,
			32626667,
			32626663,
		}

		got, err := fetcher.FetchBestStories()
		if err != nil {
			t.Errorf("Can't fetch best stories")
		}

		// Check only the first 5 ids
		for i := 0; i < 5; i++ {
			if got[i] != expected[i] {
				t.Errorf("got %d expected %d", got, expected)
			}
		}
	})

	t.Run("fetch newest stories", func(t *testing.T) {
		api := NewApi(server.Client(), server.URL)
		store := NewModelsMock()
		config := config{fetchInterval: time.Second * 5}
		logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

		fetcher := NewFetcher(config, logger, api, store)
		expected := []int{
			32411232,
			32627286,
			32626745,
			32626667,
			32626663,
		}

		got, err := fetcher.FetchNewestStories()
		if err != nil {
			t.Errorf("Can't fetch newest stories")
		}

		// Check only the first 5 ids
		for i := 0; i < 5; i++ {
			if got[i] != expected[i] {
				t.Errorf("got %d expected %d", got, expected)
			}
		}
	})

	t.Run("fetch item of type story", func(t *testing.T) {
		api := NewApi(server.Client(), server.URL)
		store := NewModelsMock()
		config := config{fetchInterval: time.Second * 5}
		logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

		fetcher := NewFetcher(config, logger, api, store)
		expected := data.Item{
			Id:   32411232,
			Type: "story",
			Kids: &[]int{
				32411282,
			},
		}

		got, err := fetcher.FetchItem(32411232)
		if err != nil {
			t.Errorf("Can't fetch newest stories")
		}

		if !reflect.DeepEqual(*got, expected) {
			t.Errorf("got %v expected %v", got, expected)
		}
	})

	t.Run("fetch comments of a story", func(t *testing.T) {
		api := NewApi(server.Client(), server.URL)
		store := NewModelsMock()
		config := config{fetchInterval: time.Second * 5}
		logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

		fetcher := NewFetcher(config, logger, api, store)

		tt := []struct {
			story    data.Item
			expected []data.Item
		}{
			{
				story: data.Item{
					Id:   32411232,
					Type: "story",
					Kids: &[]int{
						32411282,
					},
				},
				expected: []data.Item{{
					Id:   32411282,
					Type: "comment",
				}},
			},
			{
				story: data.Item{
					Id:   32627286,
					Type: "story",
					Kids: &[]int{
						32627419,
						32627299,
					},
				},
				expected: []data.Item{
					{
						Id:   32627419,
						Type: "comment",
					},
					{
						Id:   32627299,
						Type: "comment",
					},
				},
			},
			{
				story: data.Item{
					Id:   32626745,
					Type: "story",
					Kids: &[]int{
						32627477,
						32626746,
					},
				},
				expected: []data.Item{
					{
						Id:   32627477,
						Type: "comment",
					},
					{
						Id:   32626746,
						Type: "comment",
					},
				},
			},
			{
				story: data.Item{
					Id:   32626746,
					Type: "story",
					Kids: &[]int{
						32626668,
						32626685,
					},
				},
				expected: []data.Item{
					{
						Id:   32626668,
						Type: "comment",
					},
					{
						Id:   32626685,
						Type: "comment",
					},
				},
			},
			{
				story: data.Item{
					Id:   32626663,
					Type: "story",
					Kids: &[]int{
						32626998,
						32627287,
						32626728,
					},
				},
				expected: []data.Item{
					{
						Id:   32626998,
						Type: "comment",
						Kids: &[]int{
							32627096,
							32627156,
						},
					},
					{
						Id:   32627096,
						Type: "comment",
					},
					{
						Id:   32627287,
						Type: "comment",
					},
					{
						Id:   32626728,
						Type: "comment",
					},
					{
						Id:   32627156,
						Type: "comment",
					},
				},
			},
		}

		for _, testCase := range tt {
			got, err := fetcher.FetchComments(&testCase.story)
			if err != nil {
				t.Errorf("Can't fetch comments for story %v. %s", &testCase.story, err)
			}

			for _, v := range *got {
				if !Includes(testCase.expected, func(item data.Item) bool {
					return v.Id == item.Id
				}) {
					t.Errorf("got %v, expected %v", *got, testCase.expected)
				}
			}
		}

	})

	t.Run("start running", func(t *testing.T) {
		api := NewApi(server.Client(), server.URL)
		store := NewModelsMock()
		config := config{fetchInterval: 0}
		logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

		fetcher := NewFetcher(config, logger, api, store)
		expected := map[int][]int{
			32411232: {
				32411282,
			},
			32627286: {
				32627419,
				32627299,
			},
			32626745: {
				32627477,
				32626746,
			},
			32626667: {
				32626668,
				32626685,
			},
			32626663: {
				32626998,
				32627096,
				32627287,
				32626728,
				32627156,
			},
		}

		fetcher.Start()
		for storyId, comments := range store.state {
			if _, ok := expected[storyId]; !ok {
				t.Errorf("missing comments for story %d in state", storyId)
			}

			for _, comment := range expected[storyId] {
				if !Includes(comments, func(commentId int) bool {
					return comment == commentId
				}) {
					t.Errorf("expected %v to be in state", comment)
				}
			}
		}
	})
}
