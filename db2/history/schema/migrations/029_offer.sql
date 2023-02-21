-- +migrate Up

CREATE TABLE history_offer
(
    id                  BIGSERIAL   NOT NULL,
    offer_id            BIGINT      NOT NULL,
    owner_id            VARCHAR(56) NOT NULL,
    base_asset          VARCHAR(16) NOT NULL,
    quote_asset         VARCHAR(16) NOT NULL,
    is_buy              BOOLEAN     NOT NULL,
    initial_base_amount BIGINT      NOT NULL,
    current_base_amount BIGINT      NOT NULL,
    price               BIGINT      NOT NULL,
    is_canceled         BOOLEAN     NOT NULL,
    created_at  timestamp without time zone NOT NULL,
    PRIMARY KEY (id)
);

-- +migrate Down
DROP TABLE IF EXISTS history_offer;