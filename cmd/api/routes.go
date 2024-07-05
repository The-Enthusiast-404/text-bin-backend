package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	// initialize a new httprouter instance
	router := httprouter.New()

	// Convert the notFoundResponse() helper to a http.Handler using the
	// http.HandlerFunc() adapter, and then set it as the custom error handler for 404
	// Not Found responses.
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	// Convert the methodNotAllowedResponse() helper to a http.Handler using the http.HandlerFunc() adapter, and then set it as the custom error handler for 405 Method Not Allowed responses.
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// register the healthcheck handler function with the router
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthCheckHandler)

	router.HandlerFunc(http.MethodPost, "/v1/texts", app.createTextHandler)
	router.HandlerFunc(http.MethodGet, "/v1/texts/:id", app.showTextHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/texts/:id", app.updateTextHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/texts/:id", app.deleteTextHandler)

	// return the router
	return app.recoverPanic(router)
}
