-- +goose Up

CREATE TABLE news (
    id BIGINT PRIMARY KEY,
    source BIGINT NOT NULL REFERENCES source(id),  -- Changed VARCHAR(255) to BIGINT
    author VARCHAR(255),
    title VARCHAR(1000),
    description TEXT,
    publishedAt TIMESTAMPTZ
);

-- +goose Down
DROP TABLE news;
