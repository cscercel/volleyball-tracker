-- name: CreatePlayer :one
INSERT INTO players (name)
VALUES ($1)
RETURNING *;

-- name: ListPlayers :many
SELECT * FROM players
ORDER BY name;

-- name: GetPlayerByID :one
SELECT * FROM players
WHERE id = $1;

-- name: GetPlayerByName :one
SELECT * FROM players
WHERE name = $1;

-- name: UpdatePlayerName :one
UPDATE players
SET 
    name = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeletePlayer :exec
DELETE FROM players
WHERE id = $1;

-- name: UpsertPlayerStats :one
INSERT INTO player_stats (player_id, match_type, season) 
VALUES ($1, $2, $3)
ON CONFLICT (player_id, match_type, season) DO NOTHING
RETURNING *;

-- name: GetPlayerStatsByID :one
SELECT 
    p.name,
    ps.*,
    2 * ps.wins + ps.otl AS points,
    ps.wins + ps.losses + ps.otl AS played,
    CASE 
        WHEN (ps.wins + ps.losses + ps.otl) = 0 THEN 0
        ELSE CAST(ps.wins AS FLOAT) / CAST((ps.wins + ps.losses + ps.otl) AS FLOAT)
    END AS win_rate,
    CASE 
        WHEN ps.conceded = 0 THEN 0
        ELSE CAST(ps.scored AS FLOAT) / CAST(ps.conceded AS FLOAT)
    END AS efficiency_rate
FROM player_stats ps
JOIN players p ON p.id = ps.player_id
WHERE ps.player_id = $1
AND ps.match_type = $2
AND ps.season = $3;

-- name: GetPlayerStatsByName :one
SELECT 
    p.name,
    ps.*,
    2 * ps.wins + ps.otl AS points,
    ps.wins + ps.losses + ps.otl AS played,
    CASE 
        WHEN (ps.wins + ps.losses + ps.otl) = 0 THEN 0
        ELSE CAST(ps.wins AS FLOAT) / CAST((ps.wins + ps.losses + ps.otl) AS FLOAT)
    END AS win_rate,
    CASE 
        WHEN ps.conceded = 0 THEN 0
        ELSE CAST(ps.scored AS FLOAT) / CAST(ps.conceded AS FLOAT)
    END AS efficiency_rate
FROM player_stats ps
JOIN players p ON p.id = ps.player_id
WHERE p.name = $1
AND ps.match_type = $2
AND ps.season = $3;

-- name: GetLeaderboard :many
SELECT 
    p.name,
    ps.*,
    2 * ps.wins + ps.otl AS points,
    ps.wins + ps.losses + ps.otl AS played,
    CASE 
        WHEN (ps.wins + ps.losses + ps.otl) = 0 THEN 0
        ELSE CAST(ps.wins AS FLOAT) / CAST((ps.wins + ps.losses + ps.otl) AS FLOAT)
    END AS win_rate,
    CASE 
        WHEN ps.conceded = 0 THEN 0
        ELSE CAST(ps.scored AS FLOAT) / CAST(ps.conceded AS FLOAT)
    END AS efficiency_rate
FROM player_stats ps
JOIN players p ON p.id = ps.player_id
WHERE ps.match_type = $1
AND ps.season = $2
ORDER BY points DESC;

-- name: UpdatePlayerStatsWin :one
UPDATE player_stats
SET 
    wins = wins + 1,
    streak = streak + 1,
    longest_streak = GREATEST(streak + 1, longest_streak),
    scored = scored + $4,
    conceded = conceded + $5,
    updated_at = NOW()
WHERE player_id = $1
AND match_type = $2
AND season = $3
RETURNING *;

-- name: UpdatePlayerStatsLoss :one
UPDATE player_stats
SET 
    losses = losses + 1,
    streak = 0,
    scored = scored + $4,
    conceded = conceded + $5,
    updated_at = NOW()
WHERE player_id = $1
AND match_type = $2
AND season = $3
RETURNING *;

-- name: UpdatePlayerStatsOtl :one
UPDATE player_stats
SET 
    otl = otl + 1,
    scored = scored + $4,
    conceded = conceded + $5,
    updated_at = NOW()
WHERE player_id = $1
AND match_type = $2
AND season = $3
RETURNING *;
