package main

import (
	"fmt"
	"net/http"
)

// createTextHandler will be used to create a text
func (app *application) createTextHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new text")
}

// showTextHandler will be used to show a text
func (app *application) showTextHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	// show the text with id
	fmt.Fprintf(w, "show the text with id %d\n", id)
}
