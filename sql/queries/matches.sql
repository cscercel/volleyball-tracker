-- name: CreateMatch :one
INSERT INTO matches (match_type, season)
VALUES ($1, $2)
RETURNING *;

-- name: GetMatch :one
SELECT * FROM matches
WHERE id = $1;

-- name: ListMatchesBySeason :many
SELECT * FROM matches
WHERE match_type = $1
AND season = $2
AND is_completed = TRUE
ORDER BY created_at DESC;

-- name: ListUncompletedMatches :many
SELECT * FROM matches
WHERE is_completed = FALSE;

-- name: DeleteUncompletedMatch :exec
DELETE FROM matches
WHERE is_completed = FALSE
AND id = $1;

-- name: AddPlayerToMatch :one
INSERT INTO match_players (match_id, player_id, color)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetMatchPlayers :many
SELECT 
    mp.color,
    p.id AS player_id,
    p.name AS player_name
FROM match_players mp
JOIN players p ON p.id = mp.player_id
WHERE mp.match_id = $1
ORDER BY mp.color;

-- name: GetPlayerSeasonalMatches :many
SELECT
    m.id,
    m.match_type,
    m.season,
    m.blue_score,
    m.red_score,
    m.created_at,
    mp.color
FROM match_players mp
JOIN matches m ON m.id = mp.match_id
WHERE mp.player_id = $1
AND m.match_type = $2
AND m.season = $3
AND m.is_completed = TRUE
ORDER BY m.created_at DESC;
