-- +migrate Up

create table deferred_payments
(
    id    bigint      not null,
    amount NUMERIC(20, 0) not null,
    details jsonb,
    source_account VARCHAR(56) not null,
    source_balance VARCHAR(56) not null,
    destination_account VARCHAR(56) not null,
    state int not null,
    primary key (id)
);


-- +migrate Down

drop table if exists deferred_payments cascade;