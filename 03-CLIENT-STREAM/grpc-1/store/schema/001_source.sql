-- +goose Up

CREATE TABLE source (
    id BIGINT PRIMARY KEY,  -- Changed BIGINT to VARCHAR(255)
    source_id VARCHAR(255) NOT NULL UNIQUE,
    source_name VARCHAR(255) NOT NULL
);

-- +goose Down
DROP TABLE source;
