package server

import (
	"github.com/dmitryshur/hackernews/internal/validator"
	"net/http"
)

// TODO: stories type (newest or best)
func (app *application) listStoriesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Type = app.readString(qs, "type", "best")

	if ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	stories, err := app.store.GetStories(input.Type)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"stories": stories}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
