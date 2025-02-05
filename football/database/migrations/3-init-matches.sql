CREATE TYPE status_name_enum AS ENUM (
    'STATUS_FINAL_PEN',
    'STATUS_SCHEDULED',
    'STATUS_FULL_TIME',
    'STATUS_SECOND_HALF',
    'STATUS_HALF_TIME',
    'STATUS_FIRST_HALF',
    'STATUS_UNKNOWN',
    'STATUS_POSTPONED'
);

CREATE TABLE IF NOT EXISTS matches(
    id BIGINT PRIMARY KEY,
    match_date TIMESTAMP WITH TIME ZONE NOT NULL,
    venue VARCHAR NOT NULL,
    home_team_score INT NULL,
    home_team_id BIGINT NOT NULL,
    away_team_score INT NULL,
    away_team_id BIGINT NOT NULL,
    note VARCHAR(150),
    completed BOOLEAN NOT NULL,
    status_name status_name_enum NOT NULL,
    league_id INT NOT NULL,
    events JSONB NULL,
    FOREIGN KEY(home_team_id) REFERENCES teams(id),
    FOREIGN KEY(away_team_id) REFERENCES teams(id)
    FOREIGN KEY(league_id) REFERENCES leagues(id)
);