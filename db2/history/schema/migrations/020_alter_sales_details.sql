-- +migrate Up
ALTER TABLE sale ALTER COLUMN details TYPE jsonb USING to_jsonb(details::json);

-- +migrate Down
ALTER TABLE sale ALTER COLUMN details TYPE text;
