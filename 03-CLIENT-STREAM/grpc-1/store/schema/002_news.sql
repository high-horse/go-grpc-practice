-- +goose Up

CREATE TABLE news (
    id BIGSERIAL PRIMARY KEY,  -- Primary key for the news table
    source VARCHAR(255) NOT NULL REFERENCES source(source_id),  -- Foreign key reference to the source table
    author VARCHAR(255),  -- Author of the news
    title VARCHAR(1000) UNIQUE,  -- Title of the news, must be unique if required
    description TEXT,  -- Description of the news
    url VARCHAR(1000),
    publishedAt TIMESTAMPTZ  -- Timestamp when the news was published
);

-- +goose Down
DROP TABLE news;
