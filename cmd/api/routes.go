package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	// initialize a new httprouter instance
	router := httprouter.New()

	// register the healthcheck handler function with the router
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthCheckHandler)

	router.HandlerFunc(http.MethodPost, "/v1/texts", app.createTextHandler)
	router.HandlerFunc(http.MethodGet, "/v1/texts/:id", app.showTextHandler)

	// return the router
	return router
}
