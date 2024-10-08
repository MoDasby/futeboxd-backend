CREATE TABLE IF NOT EXISTS matches(
    id BIGINT PRIMARY KEY,
    match_date DATE NOT NULL,
    venue JSONB NOT NULL,
    competitors JSONB NOT NULL,
    note VARCHAR(100),
    completed BOOLEAN NOT NULL,
    status_name VARCHAR(40) NOT NULL,
    competition_name VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS teams(
    id BIGINT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    abbreviation VARCHAR(20) NOT NULL,
    color VARCHAR(6) NOT NULL DEFAULT '000000',
    logo VARCHAR(255)
);