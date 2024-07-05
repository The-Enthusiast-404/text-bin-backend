package main

import (
	"fmt"
	"net/http"
	"time"

	"dev.theenthusiast.text-bin/internal/data"
	"dev.theenthusiast.text-bin/internal/validator"
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
	err := app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	text := &data.Text{
		Title:   input.Title,
		Content: input.Content,
		Format:  input.Format,
	}

	// initialize a new validator instance
	v := validator.New()

	if data.ValidateText(v, text); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Texts.Insert(text)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/texts/%d", text.ID))
	err = app.writeJSON(w, http.StatusCreated, envelope{"text": text}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
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
