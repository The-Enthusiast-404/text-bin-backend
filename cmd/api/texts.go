package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"dev.theenthusiast.text-bin/internal/data"
	"dev.theenthusiast.text-bin/internal/validator"
)

// createTextHandler will be used to create a text
func (app *application) createTextHandler(w http.ResponseWriter, r *http.Request) {
	// Declare an anonymous struct to hold the input data that we expect to get from the request body.
	var input struct {
		Title        string `json:"title"`
		Content      string `json:"content"`
		Format       string `json:"format"`
		ExpiresValue int    `json:"expiresValue"`
		ExpiresUnit  string `json:"expiresUnit"`
	}

	// Decode the JSON request body into the input struct.
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Calculate the expiration time based on the ExpiresValue and ExpiresUnit
	var expires time.Time
	switch input.ExpiresUnit {
	case "seconds":
		expires = time.Now().Add(time.Duration(input.ExpiresValue) * time.Second)
	case "minutes":
		expires = time.Now().Add(time.Duration(input.ExpiresValue) * time.Minute)
	case "hours":
		expires = time.Now().Add(time.Duration(input.ExpiresValue) * time.Hour)
	case "days":
		expires = time.Now().Add(time.Duration(input.ExpiresValue) * time.Hour * 24)
	case "weeks":
		expires = time.Now().Add(time.Duration(input.ExpiresValue) * time.Hour * 24 * 7)
	case "months":
		expires = time.Now().AddDate(0, input.ExpiresValue, 0)
	case "years":
		expires = time.Now().AddDate(input.ExpiresValue, 0, 0)
	default:
		app.badRequestResponse(w, r, fmt.Errorf("invalid expires unit: %v", input.ExpiresUnit))
		return
	}

	// Create a new Text struct and populate it with the input data.
	text := &data.Text{
		Title:   input.Title,
		Content: input.Content,
		Format:  input.Format,
		Expires: expires,
	}

	// Initialize a new validator instance.
	v := validator.New()

	// Validate the text struct.
	if data.ValidateText(v, text); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Insert the text into the database.
	err = app.models.Texts.Insert(text)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Set the Location header for the newly created resource.
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/texts/%d", text.ID))

	// Write the JSON response with the created text.
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
		Title        *string `json:"title"`
		Content      *string `json:"content"`
		Format       *string `json:"format"`
		ExpiresUnit  *string `json:"expiresUnit"`
		ExpiresValue *int    `json:"expiresValue"`
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
	if input.ExpiresUnit != nil && input.ExpiresValue != nil {
		switch *input.ExpiresUnit {
		case "seconds":
			text.Expires = time.Now().Add(time.Duration(*input.ExpiresValue) * time.Second)
		case "minutes":
			text.Expires = time.Now().Add(time.Duration(*input.ExpiresValue) * time.Minute)
		case "hours":
			text.Expires = time.Now().Add(time.Duration(*input.ExpiresValue) * time.Hour)
		case "days":
			text.Expires = time.Now().Add(time.Duration(*input.ExpiresValue) * time.Hour * 24)
		case "weeks":
			text.Expires = time.Now().Add(time.Duration(*input.ExpiresValue) * time.Hour * 24 * 7)
		case "months":
			text.Expires = time.Now().AddDate(0, *input.ExpiresValue, 0)
		case "years":
			text.Expires = time.Now().AddDate(*input.ExpiresValue, 0, 0)
		default:
			app.badRequestResponse(w, r, fmt.Errorf("invalid expires unit: %v", *input.ExpiresUnit))
			return
		}
	}

	// Initialize a new validator instance and validate the text
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
