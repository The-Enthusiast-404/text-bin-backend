package data

import (
	"context"
	"database/sql"
	"time"
)

type Like struct {
    ID        int64     `json:"id"`
    UserID    int64     `json:"user_id"`
    TextID    int64     `json:"text_id"`
    CreatedAt time.Time `json:"created_at"`
}

type LikeModel struct {
    DB *sql.DB
}

func (m LikeModel) AddLike(userID, textID int64) error {
    query := `
        INSERT INTO likes (user_id, text_id)
        VALUES ($1, $2)
        ON CONFLICT (user_id, text_id) DO NOTHING`

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    _, err := m.DB.ExecContext(ctx, query, userID, textID)
    return err
}

func (m LikeModel) RemoveLike(userID, textID int64) error {
    query := `
        DELETE FROM likes
        WHERE user_id = $1 AND text_id = $2`

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    _, err := m.DB.ExecContext(ctx, query, userID, textID)
    return err
}
