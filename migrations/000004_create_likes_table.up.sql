CREATE TABLE IF NOT EXISTS likes (
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    text_id bigint NOT NULL REFERENCES texts ON DELETE CASCADE,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, text_id)
);
