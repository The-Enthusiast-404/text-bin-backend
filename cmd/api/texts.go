package main

import (
	"fmt"
	"net/http"
	"time"

	"dev.theenthusiast.text-bin/internal/data"
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
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
