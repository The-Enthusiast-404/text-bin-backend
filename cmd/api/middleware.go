package main

import (
	"expvar"
	"net/http"
	"strconv"

	"github.com/felixge/httpsnoop"
)

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

func (app *application) metrics(next http.Handler) http.Handler {
	totalRequestsReceived := expvar.NewInt("total_requests_received")
	totalResponsesSent := expvar.NewInt("total_responses_sent")
	totalProcessingTimeMicroseconds := expvar.NewInt("total_processing_time_μs")
	// Declare a new expvar map to hold the count of responses for each HTTP status
	// code.
	totalResponsesSentByStatus := expvar.NewMap("total_responses_sent_by_status")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Increment the requests received count, like before.
		totalRequestsReceived.Add(1)

		// Call the httpsnoop.CaptureMetrics() function, passing in the next handler in
		// the chain along with the existing http.ResponseWriter and http.Request. This
		// returns the metrics struct that we saw above.
		metrics := httpsnoop.CaptureMetrics(next, w, r)

		// Increment the response sent count, like before.
		totalResponsesSent.Add(1)

		// Get the request processing time in microseconds from httpsnoop and increment
		// the cumulative processing time.
		totalProcessingTimeMicroseconds.Add(metrics.Duration.Microseconds())

		// Use the Add() method to increment the count for the given status code by 1.
		// Note that the expvar map is string-keyed, so we need to use the strconv.Itoa()
		// function to convert the status code (which is an integer) to a string.
		totalResponsesSentByStatus.Add(strconv.Itoa(metrics.Code), 1)
	})
}
