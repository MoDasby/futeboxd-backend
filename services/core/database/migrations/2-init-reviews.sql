CREATE TABLE IF NOT EXISTS reviews(
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    rate INT CHECK (rate BETWEEN 0 AND 5) NOT NULL,
    description TEXT NOT NULL,
    match_id BIGINT NOT NULL,
    home_team_id BIGINT NOT NULL,
    away_team_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (user_id) REFERENCES users(id)
);