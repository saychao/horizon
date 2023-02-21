-- +migrate Up
ALTER TABLE history_balances DROP COLUMN kyc;

-- +migrate Down
ALTER TABLE history_balances ADD COLUMN kyc jsonb;