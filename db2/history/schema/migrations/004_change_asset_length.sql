-- +migrate Up
ALTER TABLE history_balances ALTER COLUMN asset TYPE character varying(4);

-- +migrate Down
ALTER TABLE history_balances ALTER COLUMN asset TYPE character varying(3) USING substr(asset, 1, 3);