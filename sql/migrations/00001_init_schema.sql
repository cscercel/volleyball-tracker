-- +goose Up
CREATE TABLE users (
    id              INTEGER     PRIMARY KEY AUTOINCREMENT,
    email           TEXT        NOT NULL UNIQUE,
    hashed_password TEXT        NOT NULL,
    created_at      TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE players (
    id          INTEGER     PRIMARY KEY AUTOINCREMENT,
    name        TEXT        NOT NULL UNIQUE,
    created_at  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE player_stats (
    id              INTEGER     PRIMARY KEY AUTOINCREMENT,
    player_id       INTEGER     NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    match_type      TEXT        NOT NULL CHECK (match_type IN ('indoor', 'beach')),
    season          INTEGER     NOT NULL,
    wins            INTEGER     NOT NULL DEFAULT 0,
    losses          INTEGER     NOT NULL DEFAULT 0,
    otl             INTEGER     NOT NULL DEFAULT 0,
    streak          INTEGER     NOT NULL DEFAULT 0,
    longest_streak  INTEGER     NOT NULL DEFAULT 0,
    scored          INTEGER     NOT NULL DEFAULT 0,
    conceded        INTEGER     NOT NULL DEFAULT 0,
    created_at      TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(player_id, match_type, season)
);

CREATE TABLE matches (
    id              INTEGER     PRIMARY KEY AUTOINCREMENT,
    match_type      TEXT        NOT NULL CHECK (match_type IN ('indoor', 'beach')),
    season          INTEGER     NOT NULL,
    blue_score      INTEGER     NOT NULL DEFAULT 0,
    red_score       INTEGER     NOT NULL DEFAULT 0,
    is_completed    BOOLEAN     NOT NULL DEFAULT 0,
    created_at      TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE match_players (
    id          INTEGER     PRIMARY KEY AUTOINCREMENT,
    match_id    INTEGER     NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    player_id   INTEGER     NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    color       TEXT        NOT NULL CHECK (color IN ('blue', 'red')),
    created_at  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_match_players_match_id   ON match_players(match_id);
CREATE INDEX idx_match_players_player_id  ON match_players(player_id);
CREATE INDEX idx_matches_type_season      ON matches(match_type, season);
CREATE INDEX idx_player_stats_type_season ON player_stats(match_type, season);

-- +goose Down
DROP TABLE match_players;
DROP TABLE player_stats;
DROP TABLE matches;
DROP TABLE players;
DROP TABLE users;
