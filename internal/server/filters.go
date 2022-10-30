package server

import (
	"github.com/dmitryshur/hackernews/internal/validator"
)

type ListStoriesFilters struct {
	Type string
}

type ListCommentsFilters struct {
	StoryId int64
}

func ValidateListStoriesFilters(v *validator.Validator, f ListStoriesFilters) {
	v.Check(f.Type == "newest" || f.Type == "best", "type", "must be 'newest' or 'best'")
}

func ValidateListCommentsFilters(v *validator.Validator, f ListCommentsFilters) {
	v.Check(f.StoryId > 0, "storyId", "must be higher than 0")
}
