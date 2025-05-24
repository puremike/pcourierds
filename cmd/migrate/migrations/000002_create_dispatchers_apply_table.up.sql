CREATE TABLE IF NOT EXISTS dispatchers_apply (
    id BIGSERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    vehicle TEXT NOT NULL,
    license TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);