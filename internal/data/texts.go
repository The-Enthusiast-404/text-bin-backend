package data

import (
	"time"

	"dev.theenthusiast.text-bin/internal/validator"
)

// Its important in Go to keep the Fields of a struct in Capotal letter to make it public
// Any field that starts with a lowercase letter is private to the package and aren't  exported and won't be included when encoding a struct to JSON
type Text struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Format    string    `json:"format"`
	Version   int32     `json:"version"`
}

// ValidateText will be used to validate the input data for the Text struct
func ValidateText(v *validator.Validator, text *Text) {
	v.Check(text.Title != "", "title", "must be provided")
	v.Check(len(text.Title) <= 100, "title", "must not be more than 100 bytes long")
	v.Check(text.Content != "", "content", "must be provided")
	v.Check(len(text.Content) <= 1000000, "content", "must not be more than 1000000 bytes long")
	v.Check(text.Format != "", "format", "must be provided")
}
