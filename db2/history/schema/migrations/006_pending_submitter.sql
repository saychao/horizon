-- +migrate Up

CREATE TABLE pending_transaction_signers (
    id bigserial,
    pending_transaction_id bigint NOT NULL,
    signer_identity bigint NOT NULL,
    signer_public_key character varying(64) NOT NULL
);

CREATE TABLE pending_transactions (
    id bigserial,
    tx_hash character varying(64) NOT NULL,
    tx_envelope character varying(1024) NOT NULL,
    operation_type integer NOT NULL,
    operation_key character varying(64) NOT NULL,
    state integer NOT NULL,
    tx_result character varying(255) NOT NULL,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    initiator character varying(64) NOT NULL
);

-- +migrate Down

drop table pending_transaction_signers;
drop table pending_transactions;
