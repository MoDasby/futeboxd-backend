CREATE TABLE IF NOT EXISTS matches(
    id BIGINT PRIMARY KEY,
    match_date DATE NOT NULL,
    venue VARCHAR NOT NULL,
    home_team_score INT NULL,
    home_team_id BIGINT NOT NULL,
    home_team_rosters JSONB NULL,
    away_team_score INT NULL,
    away_team_id BIGINT NOT NULL,
    away_team_rosters JSONB NULL,
    note VARCHAR(100),
    completed BOOLEAN NOT NULL,
    status_name VARCHAR(40) NOT NULL,
    competition_name VARCHAR(100) NULL, -- TODO lembrar de colocar foreign key no banco
    events JSONB NULL,
    FOREIGN KEY(home_team_id) REFERENCES teams(id),
    FOREIGN KEY(away_team_id) REFERENCES teams(id)
);

CREATE TABLE IF NOT EXISTS teams(
    id BIGINT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    abbreviation VARCHAR(20) NOT NULL,
    color VARCHAR(6) NOT NULL DEFAULT '000000',
    logo VARCHAR(255)
);