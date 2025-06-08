ALTER TABLE users
ADD COLUMN name VARCHAR(50) UNIQUE NOT NULL DEFAULT '';

-- Убираем default после добавления колонки
ALTER TABLE users ALTER COLUMN name DROP DEFAULT;

CREATE INDEX idx_users_name ON users(name);