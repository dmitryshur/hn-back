package server

import (
	"github.com/dmitryshur/hackernews/internal/validator"
	"net/http"
)

func (app *application) listCommentsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ListCommentsFilters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.StoryId = app.readInt(qs, "storyId", 0, v)

	if ValidateListCommentsFilters(v, input.ListCommentsFilters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	comments, err := app.store.GetComments(input.StoryId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"comments": comments}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
