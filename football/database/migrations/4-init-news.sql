CREATE TABLE IF NOT EXISTS news(
    id SERIAL NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    link TEXT UNIQUE NOT NULL,
    image_link TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp
)