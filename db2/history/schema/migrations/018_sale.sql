-- +migrate Up

CREATE TABLE sale
(
  id           BIGINT        NOT NULL CHECK (id >= 0),
  owner_id     VARCHAR(56)   NOT NULL,
  base_asset   VARCHAR(16)   NOT NULL,
  quote_asset  VARCHAR(16)   NOT NULL,
  start_time   TIMESTAMP without time zone NOT NULL,
  end_time     TIMESTAMP without time zone NOT NULL,
  price        NUMERIC(20,0) NOT NULL CHECK (price > 0),
  soft_cap     NUMERIC(20,0) NOT NULL CHECK (soft_cap >= 0),
  hard_cap     NUMERIC(20,0) NOT NULL CHECK (hard_cap >= 0),
  current_cap  NUMERIC(20,0) NOT NULL CHECK (current_cap >= 0),
  details      TEXT          NOT NULL,
  state        INT           NOT NULL,
  PRIMARY KEY (id)
);

-- +migrate Down

DROP TABLE IF EXISTS sale;
