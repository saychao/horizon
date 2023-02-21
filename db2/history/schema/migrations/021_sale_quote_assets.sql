-- +migrate Up
ALTER TABLE sale ADD COLUMN quote_assets jsonb;
ALTER TABLE sale DROP COLUMN quote_asset;
ALTER TABLE sale DROP COLUMN price;
ALTER TABLE sale ADD COLUMN default_quote_asset VARCHAR(16)   NOT NULL;

-- +migrate Down
ALTER TABLE sale DROP COLUMN quote_assets;
ALTER TABLE sale DROP COLUMN default_quote_asset;
ALTER TABLE sale ADD COLUMN  quote_asset VARCHAR(16);
ALTER TABLE sale ADD COLUMN price NUMERIC(20,0) CHECK (price > 0);
