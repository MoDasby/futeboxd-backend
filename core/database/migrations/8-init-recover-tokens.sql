CREATE TABLE IF NOT EXISTS recover_password_tokens(
    token VARCHAR(100) UNIQUE NOT NULL,
    user_id UUID NOT NULL,
    expires_at TIMESTAMP WITH TIMEZONE DEFAULT now() + INTERVAL '5 minutes',
    FOREIGN KEY (user_id) REFERENCES users(id)
);