-- name: CreateSource :one
INSERT INTO source
(
    source_id,
    source_name
) VALUES (
    $1, $2
) ON CONFLICT (source_id) DO UPDATE
    SET source_name = EXCLUDED.source_name
RETURNING *;

-- name: GetSourceById :one
SELECT * FROM source where id = $1 LIMIT 1;

-- name: GetSourceBySourceId :one
SELECT *  FROM source where source_id = $1 LIMIT 1;