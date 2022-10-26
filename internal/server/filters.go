package server

import (
	"github.com/dmitryshur/hackernews/internal/validator"
)

type Filters struct {
	Type string
}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Type == "newest" || f.Type == "best", "type", "must be 'newest' or 'best'")
}
