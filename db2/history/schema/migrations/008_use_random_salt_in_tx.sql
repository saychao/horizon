-- +migrate Up
ALTER TABLE history_transactions RENAME COLUMN account_sequence TO salt;

-- +migrate Down

ALTER TABLE history_transactions RENAME COLUMN salt TO account_sequence;
