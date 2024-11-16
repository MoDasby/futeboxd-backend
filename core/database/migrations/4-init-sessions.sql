CREATE TABLE IF NOT EXISTS sessions(
    id SERIAL PRIMARY KEY,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
    token VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);