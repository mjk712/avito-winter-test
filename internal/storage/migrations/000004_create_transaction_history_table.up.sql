CREATE TYPE transaction_type AS ENUM ('purchase','transfer');

CREATE TABLE IF NOT EXISTS transaction_history (
    id SERIAL PRIMARY KEY,
    from_user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    to_user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    amount INTEGER NOT NULL,
    transaction_type transaction_type NOT NULL,
    merch_id INTEGER REFERENCES merch(id) ON DELETE SET NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_transaction_history_user ON transaction_history(to_user_id,from_user_id);
CREATE INDEX idx_transaction_history_timestamp ON transaction_history(timestamp);