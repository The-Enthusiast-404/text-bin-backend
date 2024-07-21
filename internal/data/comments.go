package data

import (
	"context"
	"database/sql"
	"time"
)

type Comment struct {
    ID        int64     `json:"id"`
    UserID    int64     `json:"user_id"`
    TextID    int64     `json:"text_id"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type CommentModel struct {
    DB *sql.DB
}

func (m CommentModel) AddComment(comment *Comment) error {
    query := `
        INSERT INTO comments (user_id, text_id, content)
        VALUES ($1, $2, $3)
        RETURNING id, created_at, updated_at`

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    err := m.DB.QueryRowContext(ctx, query, comment.UserID, comment.TextID, comment.Content).Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt)
    return err
}

func (m CommentModel) DeleteComment(commentID, userID int64) error {
    query := `
        DELETE FROM comments
        WHERE id = $1 AND user_id = $2`

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    result, err := m.DB.ExecContext(ctx, query, commentID, userID)
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
