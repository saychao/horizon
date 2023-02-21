-- +migrate Up
ALTER TABLE sale DROP COLUMN current_cap;

-- +migrate Down
ALTER TABLE sale ADD COLUMN current_cap NUMERIC(20,0) NOT NULL CHECK (current_cap >= 0) DEFAULT 0;
