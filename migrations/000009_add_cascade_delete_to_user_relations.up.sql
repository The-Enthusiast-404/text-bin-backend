-- Add CASCADE DELETE to texts table
ALTER TABLE texts
DROP CONSTRAINT IF EXISTS texts_user_id_fkey,
ADD CONSTRAINT texts_user_id_fkey
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Add CASCADE DELETE to comments table
ALTER TABLE comments
DROP CONSTRAINT IF EXISTS comments_user_id_fkey,
ADD CONSTRAINT comments_user_id_fkey
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Add CASCADE DELETE to likes table
ALTER TABLE likes
DROP CONSTRAINT IF EXISTS likes_user_id_fkey,
ADD CONSTRAINT likes_user_id_fkey
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Add CASCADE DELETE to tokens table
ALTER TABLE tokens
DROP CONSTRAINT IF EXISTS tokens_user_id_fkey,
ADD CONSTRAINT tokens_user_id_fkey
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
