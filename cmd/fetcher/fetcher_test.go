package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
	"time"
)

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
		store := NewStore()
		fetcher := NewFetcher(time.Second*5, api, store)
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
		store := NewStore()
		fetcher := NewFetcher(time.Second*5, api, store)
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
		store := NewStore()
		fetcher := NewFetcher(time.Second*5, api, store)
		expected := Item{
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
		store := NewStore()
		fetcher := NewFetcher(time.Second*5, api, store)

		tt := []struct {
			story    Item
			expected []Item
		}{
			{
				story: Item{
					Id:   32411232,
					Type: "story",
					Kids: &[]int{
						32411282,
					},
				},
				expected: []Item{{
					Id:   32411282,
					Type: "comment",
				}},
			},
			{
				story: Item{
					Id:   32627286,
					Type: "story",
					Kids: &[]int{
						32627419,
						32627299,
					},
				},
				expected: []Item{
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
				story: Item{
					Id:   32626745,
					Type: "story",
					Kids: &[]int{
						32627477,
						32626746,
					},
				},
				expected: []Item{
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
				story: Item{
					Id:   32626746,
					Type: "story",
					Kids: &[]int{
						32626668,
						32626685,
					},
				},
				expected: []Item{
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
				story: Item{
					Id:   32626663,
					Type: "story",
					Kids: &[]int{
						32626998,
						32627287,
						32626728,
					},
				},
				expected: []Item{
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
				if !Includes(testCase.expected, func(item Item) bool {
					return v.Id == item.Id
				}) {
					t.Errorf("got %v, expected %v", *got, testCase.expected)
				}
			}
		}

	})
}
