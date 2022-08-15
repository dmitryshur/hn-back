package fetcher

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
		case "/beststories.json":
			w.Write(mocker.BestStories())
		case "/newstories.json":
			w.Write(mocker.NewStories())
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
		fetcher := NewFetcher(time.Second*5, api)
		expected := []int{
			32384653,
			32399949,
			32399238,
			32390526,
			32384613,
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
		fetcher := NewFetcher(time.Second*5, api)
		expected := []int{
			32411232,
			32411206,
			32411179,
			32411178,
			32411174,
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
		fetcher := NewFetcher(time.Second*5, api)
		expected := Item{
			Id:   8863,
			Type: "story",
			By:   ToPointer("dhouston"),
			Time: ToPointer(1175714200),
			Kids: &[]int{
				9224,
				8917,
				8952,
				8958,
				8884,
				8887,
				8869,
				8873,
				8940,
				8908,
				9005,
				9671,
				9067,
				9055,
				8865,
				8881,
				8872,
				8955,
				10403,
				8903,
				8928,
				9125,
				8998,
				8901,
				8902,
				8907,
				8894,
				8870,
				8878,
				8980,
				8934,
				8943,
				8876,
			},
			Url:         ToPointer("http://www.getdropbox.com/u/2/screencast.html"),
			Score:       ToPointer(104),
			Title:       ToPointer("My YC app: Dropbox - Throw away your USB drive"),
			Descendants: ToPointer(71),
		}

		got, err := fetcher.FetchItem(8863)
		if err != nil {
			t.Errorf("Can't fetch newest stories")
		}

		if !reflect.DeepEqual(*got, expected) {
			t.Errorf("got %v expected %v", got, expected)
		}
	})
}
