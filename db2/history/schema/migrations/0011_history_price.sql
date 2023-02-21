-- +migrate Up

CREATE TABLE history_price (
    base_asset character varying(16) NOT NULL,
    quote_asset character varying(16) NOT NULL,
    timestamp timestamp without time zone NOT NULL,
    price double precision NOT NULL,
    PRIMARY KEY (base_asset, quote_asset, timestamp)
);

-- +migrate Down
DROP TABLE history_price;