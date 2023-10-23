-- +goose Up
CREATE TABLE IF NOT EXISTS metrics
(
    query_time TIMESTAMP,
    user_id    CHAR(36),
    query_data JSONB
);

-- +goose Down
DROP TABLE IF EXISTS metrics;