package hackernews

import (
	"net/http"
)

func (app *application) listItemsHandler(w http.ResponseWriter, r *http.Request) {
	err := app.writeJSON(w, http.StatusOK, envelope{"item": "123"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
