ALTER TABLE texts ADD COLUMN user_id bigint REFERENCES users(id);
