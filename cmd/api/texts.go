package main

import (
	"errors"
	"fmt"
	"net/http"

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

	text, err := app.models.Texts.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	//passing text envelope instead of text struct
	err = app.writeJSON(w, http.StatusOK, envelope{"text": text}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateTextHandler(w http.ResponseWriter, r *http.Request) {

	// Read the text id parameter from the URL
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Fetch the existing text record from the database
	text, err := app.models.Texts.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Declare an input struct to hold the expected data from the request body
	var input struct {
		Title   *string `json:"title"`
		Content *string `json:"content"`
		Format  *string `json:"format"`
	}

	// Read the JSON data from the request body and store it in the input struct
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		text.Title = *input.Title
	}
	if input.Content != nil {
		text.Content = *input.Content
	}
	if input.Format != nil {
		text.Format = *input.Format
	}

	v := validator.New()
	if data.ValidateText(v, text); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Update the text record in the database
	err = app.models.Texts.Update(text)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Return a JSON response containing the updated text record
	err = app.writeJSON(w, http.StatusOK, envelope{"text": text}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteTextHandler(w http.ResponseWriter, r *http.Request) {
	// Read the text id parameter from the URL
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Delete the text record from the database
	err = app.models.Texts.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)

		}
		return
	}

	// Return a 200 OK response
	err = app.writeJSON(w, http.StatusOK, nil, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
