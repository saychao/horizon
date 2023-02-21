-- +migrate Up

CREATE TABLE history_trades (
    id BIGSERIAL NOT NULL,
    order_book_id BIGINT NOT NULL,
    base_asset character varying(16) NOT NULL,
    quote_asset character varying(16) NOT NULL,
    base_amount BIGINT NOT NULL,
    quote_amount BIGINT NOT NULL,
    price BIGINT NOT NULL,
    created_at timestamp without time zone,
    PRIMARY KEY (id, base_asset, quote_asset)
);

-- +migrate Down
DROP TABLE history_trades;