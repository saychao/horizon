-- +migrate Up
ALTER TABLE history_contracts ADD COLUMN customer_details jsonb NOT NULL DEFAULT '{}';

-- +migrate Down
ALTER TABLE history_contracts DROP COLUMN IF EXISTS customer_details;
