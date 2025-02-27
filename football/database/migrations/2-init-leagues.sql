CREATE TABLE IF NOT EXISTS leagues (
    id SERIAL NOT NULL,
    name TEXT NOT NULL,
    logo TEXT NOT NULL,
    espn_id TEXT NOT NULL,
    PRIMARY KEY (id)
);