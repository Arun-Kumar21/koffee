-- +goose up
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(30) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    avatar_url TEXT,

    created_at TIMESTAMP DEFAULT NOW(),
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- +goose down
DROP TABLE IF EXISTS users;
