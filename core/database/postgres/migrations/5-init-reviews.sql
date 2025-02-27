CREATE TABLE IF NOT EXISTS reviews(
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    rate INT CHECK (rate BETWEEN 0 AND 5) NOT NULL,
    description VARCHAR(1500),
    match_id BIGINT NOT NULL,
    home_team_id BIGINT NOT NULL,
    away_team_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
    search_vector TSVECTOR,
    FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT unique_user_match UNIQUE (user_id, match_id)
);

CREATE INDEX IF NOT EXISTS idx_reviews_home_team_id ON reviews(home_team_id);
CREATE INDEX IF NOT EXISTS idx_reviews_away_team_id ON reviews(away_team_id);
CREATE INDEX IF NOT EXISTS idx_reviews_match_id ON reviews(match_id);

CREATE OR REPLACE FUNCTION update_reviews_search_vector() RETURNS trigger AS $$
BEGIN
  NEW.search_vector := to_tsvector('portuguese', unaccent(NEW.description));
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_reviews_search_vector
BEFORE INSERT OR UPDATE ON reviews
FOR EACH ROW EXECUTE FUNCTION update_reviews_search_vector();