package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"dev.theenthusiast.text-bin/internal/data"
)

// createTextHandler will be used to create a text
func (app *application) createTextHandler(w http.ResponseWriter, r *http.Request) {
	// declare a anonymous struct to hold the input data that we expect to get from the request body (the field names are subset of the Text struct)
	// This struct will be the target decode destination for the JSON decoder
	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Format  string `json:"format"`
	}
	// Initialize a new json.Decoder instance which reads from the request body, and
	// then use the Decode() method to decode the body contents into the input struct.
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Fprintf(w, "%+v\n", input)
}

// showTextHandler will be used to show a text
func (app *application) showTextHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// creating a static text instance to insert into the database
	text := data.Text{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Golang snippet",
		Content:   "This is a Golang snippet",
		Format:    "golang",
		Version:   1,
	}
	//passing text envelope instead of text struct
	err = app.writeJSON(w, http.StatusOK, envelope{"text": text}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
