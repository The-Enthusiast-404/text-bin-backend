CREATE TABLE IF NOT EXISTS texts (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT now(),
    title text NOT NULL,
    content text NOT NULL,
    format text NOT NULL DEFAULT 'plaintext',
    expires timestamp(0) with time zone,
    version integer NOT NULL DEFAULT 1
);
