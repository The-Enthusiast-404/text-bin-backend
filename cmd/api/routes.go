package main

import (
	"expvar"
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

	router.HandlerFunc(http.MethodPost, "/v1/users/email", app.getCurrentUser)
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/users/:id", app.deleteAccountHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users/authentication", app.createAuthenticationTokenHandler)

	router.HandlerFunc(http.MethodPut, "/v1/users/password", app.updateUserPasswordHandler)

	// Add the POST /v1/tokens/password-reset endpoint.
	router.HandlerFunc(http.MethodPost, "/v1/tokens/password-reset", app.createPasswordResetTokenHandler)

	router.HandlerFunc(http.MethodPost, "/v1/texts/:id/like", app.addLikeHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/texts/:id/like", app.removeLikeHandler)
	router.HandlerFunc(http.MethodPost, "/v1/texts/:id/comments", app.addCommentHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/texts/:id/comments/:commentID", app.deleteCommentHandler)

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	// return the router
	return app.enableCORS((app.metrics(app.recoverPanic(app.authenticate(router)))))
}
