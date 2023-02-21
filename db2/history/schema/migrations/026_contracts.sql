-- +migrate Up

CREATE TABLE history_contracts
(
  id              BIGINT      NOT NULL CHECK (id >= 0),
  contractor      VARCHAR(56) NOT NULL,
  customer        VARCHAR(56) NOT NULL,
  escrow          VARCHAR(56) NOT NULL,
  start_time      TIMESTAMP without time zone NOT NULL,
  end_time        TIMESTAMP without time zone NOT NULL,
  initial_details jsonb       NOT NULL,
  invoices        BIGINT[]    DEFAULT NULL,
  state           INT         NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE history_contracts_details
(
  contract_id BIGINT      NOT NULL CHECK (contract_id >= 0),
  details     jsonb       NOT NULL,
  author      VARCHAR(56) NOT NULL,
  created_at  TIMESTAMP without time zone NOT NULL
);

create index on history_contracts_details (contract_id);

CREATE TABLE history_contracts_disputes
(
  contract_id BIGINT      NOT NULL CHECK (contract_id >= 0),
  reason      jsonb       NOT NULL,
  author      VARCHAR(56) NOT NULL,
  created_at  TIMESTAMP without time zone NOT NULL
);

create index on history_contracts_disputes (contract_id);

-- +migrate Down

DROP TABLE IF EXISTS history_contracts;
DROP TABLE IF EXISTS history_contracts_details;
DROP TABLE IF EXISTS history_contracts_disputes;