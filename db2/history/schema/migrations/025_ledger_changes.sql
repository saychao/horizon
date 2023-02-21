-- +migrate Up

CREATE TABLE history_ledger_changes
(
  tx_id         BIGINT    NOT NULL,
  op_id         BIGINT    NOT NULL,
  order_number  INT       NOT NULL,
  effect        INT       NOT NULL,
  entry_type    INT       NOT NULL,
  payload       TEXT      NOT NULL
);

create index on history_ledger_changes (tx_id, entry_type);

-- +migrate Down

DROP TABLE IF EXISTS history_ledger_changes;