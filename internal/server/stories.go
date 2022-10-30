package server

import (
	"errors"
	"github.com/dmitryshur/hackernews/internal/data"
	"github.com/dmitryshur/hackernews/internal/validator"
	"net/http"
)

func (app *application) showStoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	story, err := app.store.GetStory(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"story": story}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listStoriesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ListStoriesFilters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Type = app.readString(qs, "type", "best")

	if ValidateListStoriesFilters(v, input.ListStoriesFilters); !v.Valid() {
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
