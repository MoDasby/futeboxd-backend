CREATE TABLE IF NOT EXISTS likes (
    like_owner_id VARCHAR(255) NOT NULL,
    review_id INT,
    comment_id INT,
    created_at TIMESTAMP DEFAULT current_timestamp,
    
    FOREIGN KEY (like_owner_id) REFERENCES users(id),
    FOREIGN KEY (review_id) REFERENCES reviews(id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    
    UNIQUE (like_owner_id, review_id),
    UNIQUE (like_owner_id, comment_id),
    
    CHECK (
        (review_id IS NOT NULL AND comment_id IS NULL) OR
        (review_id IS NULL AND comment_id IS NOT NULL)
    )
);