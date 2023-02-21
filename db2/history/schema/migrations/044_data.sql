-- +migrate Up

create table data
(
    id    bigint      not null,
    type  bigint      not null,
    owner varchar(56) not null,
    value jsonb       not null,
    primary key (id)
);

create index data_by_owner on data using btree (owner);
create index data_by_type on data using btree (type);

-- +migrate Down

drop table if exists data cascade;