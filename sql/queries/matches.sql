-- name: CreateMatch :one
INSERT INTO matches (match_type, season)
VALUES (?, ?)
RETURNING *;

-- name: GetMatch :one
SELECT * FROM matches
WHERE id = ?;

-- name: UpdateMatchScores :one
UPDATE matches
SET 
    blue_score = ?,
    red_score = ?,
    is_completed = 1,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: ListMatchesBySeason :many
SELECT * FROM matches
WHERE match_type = ?
AND season = ?
AND is_completed = 1
ORDER BY created_at DESC
LIMIT ?
OFFSET ?;

-- name: ListUncompletedMatches :many
SELECT * FROM matches
WHERE is_completed = 0;

-- name: DeleteUncompletedMatch :exec
DELETE FROM matches
WHERE is_completed = 0
AND id = ?;

-- name: AddPlayerToMatch :one
INSERT INTO match_players (match_id, player_id, color)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetMatchPlayers :many
SELECT 
    mp.color,
    p.id AS player_id,
    p.name AS player_name
FROM match_players mp
JOIN players p ON p.id = mp.player_id
WHERE mp.match_id = ?
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
WHERE mp.player_id = ?
AND m.match_type = ?
AND m.season = ?
AND m.is_completed = 1
ORDER BY m.created_at DESC
LIMIT ?
OFFSET ?;
