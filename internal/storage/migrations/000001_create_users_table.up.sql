CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE  NOT NULL,
    password VARCHAR(255) NOT NULL,
    coins INTEGER DEFAULT 1000
);

CREATE INDEX idx_users_username ON users(username);