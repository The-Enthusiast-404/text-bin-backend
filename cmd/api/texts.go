package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// createTextHandler will be used to create a text
func (app *application) createTextHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new text")
}

// showTextHandler will be used to show a text
func (app *application) showTextHandler(w http.ResponseWriter, r *http.Request) {
	// when httprouter is used, the parameters are available in the context of the request
	params := httprouter.ParamsFromContext(r.Context())
	// convert the id parameter to an integer
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	// if there is an error or id is less than 1, return a 404 Not Found response
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	// show the text with id
	fmt.Fprintf(w, "show the text with id %d\n", id)
}
