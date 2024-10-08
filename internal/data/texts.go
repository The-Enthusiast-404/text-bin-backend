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
	"golang.org/x/exp/rand"
)

// Its important in Go to keep the Fields of a struct in Capotal letter to make it public
// Any field that starts with a lowercase letter is private to the package and aren't  exported and won't be included when encoding a struct to JSON
type Text struct {
	ID             int64     `json:"id"`
	CreatedAt      time.Time `json:"-"`
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	Format         string    `json:"format"`
	Expires        time.Time `json:"expires"`
	Slug           string    `json:"slug"`
	IsPrivate      bool      `json:"is_private"`
	UserID         *int64    `json:"user_id,omitempty"`
	LikesCount     int       `json:"likes_count"`
	Comments       []Comment `json:"comments,omitempty"`
	EncryptionSalt string    `json:"encryption_salt"`
	Version        int32     `json:"-"`
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
        INSERT INTO texts (title, content, format, expires, slug, user_id, is_private, encryption_salt)
        VALUES($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, created_at, version
    `
	args := []interface{}{
		text.Title,
		text.Content,
		text.Format,
		text.Expires,
		text.Slug,
		text.UserID,
		text.IsPrivate,
		text.EncryptionSalt,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&text.ID, &text.CreatedAt, &text.Version)
	if err != nil {
		return fmt.Errorf("failed to insert text: %v", err)
	}
	return nil
}

var randomSource rand.Source
var randomGenerator *rand.Rand

func init() {
	// Convert int64 to uint64 without losing information
	seed := uint64(time.Now().UnixNano())
	source := rand.NewSource(seed)
	randomGenerator = rand.New(source)
}

func (m TextModel) GenerateUniqueSlug(title string) (string, error) {
	baseSlug := generateBaseSlug(title)

	for attempts := 0; attempts < 10; attempts++ {
		slug := baseSlug
		if attempts > 0 {
			// Add a short random string instead of a number
			randomStr := generateRandomString(3)
			slug = fmt.Sprintf("%s-%s", baseSlug, randomStr)
		}

		exists, err := m.slugExists(slug)
		if err != nil {
			return "", err
		}
		if !exists {
			return slug, nil
		}
	}

	// If all attempts fail, generate a completely random slug
	return generateRandomString(8), nil
}

func generateBaseSlug(title string) string {
	title = strings.ToLower(title)
	reg := regexp.MustCompile("[^a-z0-9]+")
	title = reg.ReplaceAllString(title, " ")
	words := strings.Fields(title)
	if len(words) > 3 {
		words = words[:3]
	}
	slug := strings.Join(words, "-")
	if len(slug) > 20 {
		slug = slug[:20]
	}
	return strings.TrimRight(slug, "-")
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[randomGenerator.Intn(len(charset))]
	}
	return string(b)
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
        SELECT id, created_at, title, content, format, expires, slug, version, user_id, is_private, encryption_salt,
               (SELECT COUNT(*) FROM likes WHERE text_id = texts.id) as likes_count
        FROM texts
        WHERE slug = $1`

	var text Text

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, slug).Scan(
		&text.ID, &text.CreatedAt, &text.Title, &text.Content, &text.Format,
		&text.Expires, &text.Slug, &text.Version, &text.UserID, &text.IsPrivate,
		&text.EncryptionSalt, &text.LikesCount)

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
        SET title = $1, content = $2, format = $3, expires = $4, is_private = $5, encryption_salt = $6, version = version + 1
        WHERE slug = $7 AND version = $8 AND (user_id = $9 OR user_id IS NULL)
        RETURNING version
    `
	args := []interface{}{
		text.Title,
		text.Content,
		text.Format,
		text.Expires,
		text.IsPrivate,
		text.EncryptionSalt,
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
