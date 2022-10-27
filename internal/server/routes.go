package server

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/stories", app.listStoriesHandler)
	router.HandlerFunc(http.MethodGet, "/v1/stories/:id", app.showStoryHandler)

	return app.recoverPanic(app.enableCORS(app.rateLimit(router)))
}
