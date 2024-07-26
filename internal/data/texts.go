package data

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"dev.theenthusiast.text-bin/internal/validator"
)

// Its important in Go to keep the Fields of a struct in Capotal letter to make it public
// Any field that starts with a lowercase letter is private to the package and aren't  exported and won't be included when encoding a struct to JSON
type Text struct {
	ID         int64     `json:"id"`
	CreatedAt  time.Time `json:"-"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Format     string    `json:"format"`
	Version    int32     `json:"version"`
	Expires    time.Time `json:"expires"`
	Slug       string    `json:"slug"`
	LikesCount int       `json:"likes_count"`
	Comments   []Comment `json:"comments,omitempty"`
	UserID     int64     `json:"user_id,omitempty"`
}

// ValidateText will be used to validate the input data for the Text struct
func ValidateText(v *validator.Validator, text *Text) {
	v.Check(text.Title != "", "title", "must be provided")
	v.Check(len(text.Title) <= 100, "title", "must not be more than 100 bytes long")
	v.Check(text.Content != "", "content", "must be provided")
	v.Check(len(text.Content) <= 1000000, "content", "must not be more than 1000000 bytes long")
	v.Check(text.Format != "", "format", "must be provided")
	v.Check(text.Expires.After(time.Now()), "expires", "must be greater than the current time")

}

// GenerateRandomCode generates a random string of specified length
func GenerateRandomCode(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Define a MovieModel struct type which wraps a sql.DB connection pool.
type TextModel struct {
	DB *sql.DB
}

// Insert will add a new record to the texts table
func (m TextModel) Insert(text *Text) error {
	query :=
		`
        INSERT INTO texts (title, content, format, expires, slug, user_id)
        VALUES($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at, version
        `
	args := []interface{}{text.Title, text.Content, text.Format, text.Expires, text.Slug, text.UserID}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&text.ID, &text.CreatedAt, &text.Version)
	if err != nil {
		return fmt.Errorf("failed to insert text: %v", err)
	}
	return nil
}

// Get will return a specific record from the texts table based on the id
func (m TextModel) Get(slug string) (*Text, error) {
	query := `
        SELECT t.id, t.created_at, t.title, t.content, t.format, t.expires, t.slug, t.version,
               (SELECT COUNT(*) FROM likes WHERE text_id = t.id) as likes_count
        FROM texts t
        WHERE slug = $1`

	var text Text

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, slug).Scan(
		&text.ID, &text.CreatedAt, &text.Title, &text.Content, &text.Format,
		&text.Expires, &text.Slug, &text.Version, &text.LikesCount)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	// Fetch comments
	commentsQuery := `
        SELECT id, user_id, content, created_at, updated_at
        FROM comments
        WHERE text_id = $1
        ORDER BY created_at DESC`

	rows, err := m.DB.QueryContext(ctx, commentsQuery, text.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt)
		if err != nil {
			return nil, err
		}
		text.Comments = append(text.Comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &text, nil
}

// Update will update a specific record in the texts table based on the id
func (m TextModel) Update(text *Text) error {
	query :=
		`
		UPDATE texts
		SET title = $1, content = $2, format = $3,expires = $4, version = version + 1
		WHERE slug = $5 AND version = $6
		RETURNING version
	`
	args := []interface{}{text.Title, text.Content, text.Format, text.Expires, text.Slug, text.Version}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&text.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

// Delete will remove a specific record from the texts table based on the id
func (m TextModel) Delete(id string) error {
	if id == "" {
		return ErrRecordNotFound
	}
	query :=
		`
		DELETE FROM texts
		WHERE slug = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	result, err := m.DB.ExecContext(ctx, query, id)
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
