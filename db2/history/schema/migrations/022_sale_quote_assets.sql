-- +migrate Up
ALTER TABLE sale ADD COLUMN base_current_cap NUMERIC(20,0) NOT NULL DEFAULT 0;
ALTER TABLE sale ADD COLUMN base_hard_cap NUMERIC(20,0) NOT NULL DEFAULT 0;

-- +migrate Down
ALTER TABLE sale DROP COLUMN base_current_cap;
ALTER TABLE sale DROP COLUMN base_hard_cap;
