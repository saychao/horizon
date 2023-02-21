-- +migrate Up

ALTER TABLE reviewable_request ADD COLUMN created_at timestamp without time zone NOT NULL;
ALTER TABLE reviewable_request ADD COLUMN updated_at timestamp without time zone NOT NULL;

-- +migrate Down

ALTER TABLE reviewable_request DROP COLUMN created_at;
ALTER TABLE reviewable_request DROP COLUMN updated_at;
