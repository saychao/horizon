-- +migrate Up

CREATE TABLE reviewable_request (
  id            BIGINT NOT NULL,
  requestor     VARCHAR(56)   NOT NULL,
  reviewer      VARCHAR(56)   NOT NULL,
  reference     VARCHAR(64),
  reject_reason TEXT          NOT NULL,
  request_type  INT           NOT NULL,
  request_state INT           NOT NULL,
  hash          CHARACTER(64) NOT NULL,
  details          jsonb         NOT NULL,
  PRIMARY KEY (id)
);

-- +migrate Down

DROP TABLE IF EXISTS reviewable_request;
