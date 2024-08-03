-- name: CreateNews :one
INSERT INTO news (
    source,
    author,
    title,
    description,
    publishedAt
) VALUES (
    $1, $2, $3, $4, NOW()
) RETURNING *;


-- name: GetSingleNews :one
SELECT N.*, S.source_name FROM news N
INNER JOIN source S ON S.id = N.source
WHERE N.id = $1 LIMIT 1;

-- name: GetSourceBasedNews :many
SELECT N.*, S.source_name FROM news N
INNER JOIN source S ON S.id = N.source
WHERE N.source = $1;
