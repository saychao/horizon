-- +migrate Up
ALTER TABLE sale ADD COLUMN sale_type INT NOT NULL DEFAULT 1;

-- +migrate Down
ALTER TABLE sale DROP COLUMN  sale_type;
