CREATE TABLE IF NOT EXISTS inventory(
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    merch_id INTEGER REFERENCES merch(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL,
    UNIQUE (user_id, merch_id)
);

CREATE INDEX idx_inventory_user_merch ON inventory(user_id,merch_id);