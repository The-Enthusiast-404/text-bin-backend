package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
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
	UserID     *int64    `json:"user_id,omitempty"`
	IsPrivate  bool      `json:"is_private"`
}

// ValidateText will be used to validate the input data for the Text struct
func ValidateText(v *validator.Validator, text *Text) {
	v.Check(text.Title != "", "title", "must be provided")
	v.Check(len(text.Title) <= 100, "title", "must not be more than 100 bytes long")
	v.Check(text.Content != "", "content", "must be provided")
	v.Check(len(text.Content) <= 1000000, "content", "must not be more than 1000000 bytes long")
	v.Check(text.Format != "", "format", "must be provided")
	v.Check(text.Expires.After(time.Now()), "expires", "must be greater than the current time")
	v.Check(text.UserID != nil || !text.IsPrivate, "is_private", "anonymous users cannot create private texts")

}

// GenerateRandomCode generates a random string of specified length
// func GenerateRandomCode(n int) (string, error) {
// 	b := make([]byte, n)
// 	_, err := rand.Read(b)
// 	if err != nil {
// 		return "", err
// 	}
// 	return base64.URLEncoding.EncodeToString(b), nil
// }

// Define a MovieModel struct type which wraps a sql.DB connection pool.
type TextModel struct {
	DB *sql.DB
}

// Insert will add a new record to the texts table
func (m TextModel) Insert(text *Text) error {
	query := `
        INSERT INTO texts (title, content, format, expires, slug, user_id, is_private)
        VALUES($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at, version
    `
	args := []interface{}{text.Title, text.Content, text.Format, text.Expires, text.Slug, text.UserID, text.IsPrivate}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&text.ID, &text.CreatedAt, &text.Version)
	if err != nil {
		return fmt.Errorf("failed to insert text: %v", err)
	}
	return nil
}

func (m TextModel) GenerateUniqueSlug(title string) (string, error) {
	baseSlug := strings.ToLower(strings.ReplaceAll(title, " ", "-"))
	baseSlug = regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(baseSlug, "")

	for i := 0; i < 100; i++ { // Try up to 100 times
		slug := baseSlug
		if i > 0 {
			slug = fmt.Sprintf("%s-%d", baseSlug, i)
		}

		exists, err := m.slugExists(slug)
		if err != nil {
			return "", err
		}
		if !exists {
			return slug, nil
		}
	}

	return "", errors.New("unable to generate unique slug")
}

func (m TextModel) slugExists(slug string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM texts WHERE slug = $1)"
	err := m.DB.QueryRow(query, slug).Scan(&exists)
	return exists, err
}

// Get will return a specific record from the texts table based on the id
func (m TextModel) Get(slug string, userID *int64) (*Text, error) {
	query := `
        SELECT id, created_at, title, content, format, expires, slug, version, user_id, is_private,
               (SELECT COUNT(*) FROM likes WHERE text_id = texts.id) as likes_count
        FROM texts
        WHERE slug = $1`

	var text Text

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, slug).Scan(
		&text.ID, &text.CreatedAt, &text.Title, &text.Content, &text.Format,
		&text.Expires, &text.Slug, &text.Version, &text.UserID, &text.IsPrivate, &text.LikesCount)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	// Check if the text is private and the user is not the owner
	if text.IsPrivate && (userID == nil || *userID != *text.UserID) {
		return nil, ErrRecordNotFound
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
// Update will update a specific record in the texts table based on the id
func (m TextModel) Update(text *Text, userID int64) error {
	query := `
        UPDATE texts
        SET title = $1, content = $2, format = $3, expires = $4, is_private = $5, version = version + 1
        WHERE slug = $6 AND version = $7 AND (user_id = $8 OR user_id IS NULL)
        RETURNING version
    `
	args := []interface{}{
		text.Title,
		text.Content,
		text.Format,
		text.Expires,
		text.IsPrivate,
		text.Slug,
		text.Version,
		userID,
	}

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
// Delete will remove a specific record from the texts table based on the id
func (m TextModel) Delete(slug string, userID int64) error {
	if slug == "" {
		return ErrRecordNotFound
	}
	query := `
        DELETE FROM texts
        WHERE slug = $1 AND (user_id = $2 OR user_id IS NULL)
    `
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, slug, userID)
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
