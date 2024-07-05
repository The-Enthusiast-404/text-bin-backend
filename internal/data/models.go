package data

import (
	"database/sql"
	"errors"
)

// Define a custom error type for when an expected record is not found in the database.
var (
	ErrRecordNotFound = errors.New("record not found")
)

// Define a Models type which wraps the MovieModel.
type Models struct {
	Texts TextModel
}

// Define a NewModels() function which initializes the MovieModel and stores it in the Models type.
func NewModels(db *sql.DB) Models {
	return Models{
		Texts: TextModel{DB: db},
	}
}
