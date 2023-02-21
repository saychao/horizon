-- +migrate Up
ALTER TABLE pending_transaction_signers ADD COLUMN signer_name text NOT NULL DEFAULT '';

-- +migrate Down
ALTER TABLE pending_transaction_signers DROP COLUMN signer_name;