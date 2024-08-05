-- name: CreateNews :one
INSERT INTO news (
    source,
    author,
    title,
    description,
    publishedAt
) VALUES (
    (SELECT source_id FROM source WHERE source_id = $1),
    $2, $3, $4, $5
) ON CONFLICT (title) DO UPDATE
    SET author = EXCLUDED.author,
        description = EXCLUDED.description,
        publishedAt = EXCLUDED.publishedAt
RETURNING *;


-- name: GetSingleNews :one
SELECT N.*, S.source_name FROM news N
INNER JOIN source S ON S.id = N.source
WHERE N.id = $1 LIMIT 1;

-- name: GetSourceBasedNews :many
SELECT N.*, S.source_name FROM news N
INNER JOIN source S ON S.id = N.source
WHERE N.source = $1;

-- name: GetAllNews :many
SELECT N.*, S.source_name FROM news N
INNER JOIN source S ON S.source_id = N.source ;
