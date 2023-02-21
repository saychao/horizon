-- +migrate Up


create table swaps
(
    id                      bigint                      not null,
    created_at              TIMESTAMP without time zone NOT NULL,
    lock_time               TIMESTAMP without time zone NOT NULL,

    source_account          varchar(56)                 not null,
    source_balance          varchar(56)                 not null,
    destination_account     varchar(56)                 not null,
    destination_balance     varchar(56)                 not null,

    secret_hash             varchar(64)                 not null,
    secret                  varchar(64) default null,

    asset                   varchar(16)                 not null,
    amount                  NUMERIC(20, 0)              not null,
    source_fixed_fee        NUMERIC(20, 0)              not null,
    source_percent_fee      NUMERIC(20, 0)              not null,
    destination_fixed_fee   NUMERIC(20, 0)              not null,
    destination_percent_fee NUMERIC(20, 0)              not null,

    details                 jsonb                       not null,

    state                   int                         not null,
    primary key (id)
);

create index swaps_by_source on swaps using btree (source_account);
create index swaps_by_destination on swaps using btree (destination_account);

-- +migrate Down

drop table if exists swaps cascade;
