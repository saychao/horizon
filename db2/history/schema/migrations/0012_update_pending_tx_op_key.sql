-- +migrate Up
ALTER TABLE pending_transactions ALTER tx_envelope TYPE text;
ALTER TABLE pending_transactions ALTER tx_result TYPE text;
ALTER TABLE pending_transactions ALTER operation_key DROP NOT NUll;

CREATE UNIQUE INDEX pen_tx_operation_key_unique_index
  ON pending_transactions (operation_key)
  WHERE (operation_key IS NOT NULL);


-- +migrate Down
DROP INDEX pen_tx_operation_key_unique_index;

UPDATE pending_transactions SET operation_key='op_key' WHERE operation_key IS NULL;
ALTER TABLE pending_transactions ALTER operation_key SET NOT NUll;

UPDATE pending_transactions SET tx_envelope='tx_envelope', tx_result='tx_result';
ALTER TABLE pending_transactions ALTER tx_envelope TYPE varchar(1024);
ALTER TABLE pending_transactions ALTER tx_result TYPE varchar(255);
