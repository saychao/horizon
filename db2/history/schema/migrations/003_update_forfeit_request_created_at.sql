-- +migrate Up
ALTER TABLE history_forfeit_requests ALTER COLUMN created_at TYPE BIGINT USING EXTRACT(EPOCH FROM created_at);

-- +migrate Down
ALTER TABLE history_forfeit_requests DROP COLUMN created_at;
ALTER TABLE history_forfeit_requests ADD COLUMN created_at timestamp without time zone;
