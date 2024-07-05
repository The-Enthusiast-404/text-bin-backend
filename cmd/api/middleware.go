package main

import "net/http"

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event of a panic as Go unwinds the stack).
		defer func() {
			// Use the builtin recover function to check if there has been a panic.
			if err := recover(); err != nil {
				// Call the serverErrorResponse helper method to return a 500 Internal Server Error response.
				app.serverErrorResponse(w, r, err.(error))
			}
		}()

		// Call the ServeHTTP method on the next http.Handler in the chain.
		next.ServeHTTP(w, r)
	})
}
