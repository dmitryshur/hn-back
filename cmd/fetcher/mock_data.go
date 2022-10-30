package main

import (
	"encoding/json"
	"github.com/dmitryshur/hackernews/internal/data"
	"log"
)

type Mocks struct {
}

func NewMocker() Mocks {
	return Mocks{}
}

func (m Mocks) Stories() []byte {
	newStories := []int{
		32411232,
		32627286,
		32626745,
		32626667,
		32626663,
	}
	bytes, err := json.Marshal(newStories)
	if err != nil {
		log.Fatalf("Can't create mock for new stories. %s", err)
	}

	return bytes
}

func (m Mocks) Item(id int) []byte {
	items := map[int]data.Item{
		32411232: {
			Id:   32411232,
			Type: "story",
			Kids: &[]int64{
				32411282,
			},
		},
		32411282: {
			Id:   32411282,
			Type: "comment",
		},
		32627286: {
			Id:   32627286,
			Type: "story",
			Kids: &[]int64{
				32627419,
				32627299,
			},
		},
		32627419: {
			Id:   32627419,
			Type: "comment",
		},
		32627299: {
			Id:   32627299,
			Type: "comment",
		},
		32626745: {
			Id:   32626745,
			Type: "story",
			Kids: &[]int64{
				32627477,
				32626746,
			},
		},
		32627477: {
			Id:   32627477,
			Type: "comment",
		},
		32626746: {
			Id:   32626746,
			Type: "comment",
		},
		32626667: {
			Id:   32626667,
			Type: "story",
			Kids: &[]int64{
				32626668,
				32626685,
			},
		},
		32626668: {
			Id:   32626668,
			Type: "comment",
		},
		32626685: {
			Id:   32626685,
			Type: "comment",
		},
		32626663: {
			Id:   32626663,
			Type: "story",
			Kids: &[]int64{
				32626998,
				32627287,
				32626728,
			},
		},
		32626998: {
			Id:   32626998,
			Type: "comment",
			Kids: &[]int64{
				32627096,
				32627156,
			},
		},
		32627096: {
			Id:   32627096,
			Type: "comment",
		},
		32627287: {
			Id:   32627287,
			Type: "comment",
		},
		32626728: {
			Id:   32626728,
			Type: "comment",
		},
		32627156: {
			Id:   32627156,
			Type: "comment",
		},
	}

	bytes, err := json.Marshal(items[id])
	if err != nil {
		log.Fatalf("Can't create mock item with id 8863")
	}

	return bytes
}
