-- +migrate Up
DROP TABLE history_forfeit_requests;
ALTER TABLE history_payment_requests ADD COLUMN request_type int;

-- +migrate Down

CREATE TABLE history_forfeit_requests (
    id bigint NOT NULL,
    target character varying(64) NOT NULL,
    amount character varying(64) NOT NULL,
    initiated_by_user boolean NOT NULL,
    accepted boolean,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    forfeit_type integer DEFAULT 0 NOT NULL
);

ALTER TABLE history_payment_requests DROP COLUMN request_type;