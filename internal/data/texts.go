package data

import (
	"database/sql"
	"errors"
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

// Define a MovieModel struct type which wraps a sql.DB connection pool.
type TextModel struct {
	DB *sql.DB
}

// Insert will add a new record to the texts table
func (m TextModel) Insert(text *Text) error {
	query :=
		`
			INSERT INTO texts (title, content, format)
			VALUES($1, $2, $3)
			RETURNING id, created_at, version
		`
	args := []interface{}{text.Title, text.Content, text.Format}
	return m.DB.QueryRow(query, args...).Scan(&text.ID, &text.CreatedAt, &text.Version)
}

// Get will return a specific record from the texts table based on the id
func (m TextModel) Get(id int64) (*Text, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query :=
		`
		SELECT id, created_at, title, content, format, version
		FROM texts
		WHERE id = $1
	`

	// declare a text variable to hold the data from the query
	var text Text
	err := m.DB.QueryRow(query, id).Scan(&text.ID, &text.CreatedAt, &text.Title, &text.Content, &text.Format, &text.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &text, nil

}

// Update will update a specific record in the texts table based on the id
func (m TextModel) Update(text *Text) error {
	query :=
		`
		UPDATE texts
		SET title = $1, content = $2, format = $3, version = version + 1
		WHERE id = $4
		RETURNING version
	`
	args := []interface{}{text.Title, text.Content, text.Format, text.ID}
	return m.DB.QueryRow(query, args...).Scan(&text.Version)
}

// Delete will remove a specific record from the texts table based on the id
func (m TextModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query :=
		`
		DELETE FROM texts
		WHERE id = $1
	`
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
