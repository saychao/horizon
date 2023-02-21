-- +migrate Up

create index participant_effects_sale_participations_idx on participant_effects using btree (((effect#>>'{matched,order_book_id}')::int), asset_code, ((effect#>>'{matched,offer_id}')::int)) where (effect#>>'{type}')::int = 8 and (effect#>>'{matched,order_book_id}')::int != 0;
drop table if exists sale_participation cascade;

-- +migrate Down

create table sale_participation (
    id bigint not null,
    participant_id varchar(56) not null,
    sale_id bigint not null,
    base_amount numeric(20,0) not null,
    quote_amount numeric(20,0) not null,
    quote_asset varchar(16) not null,
    base_asset varchar(16) not null,
    price numeric(20,0) not null,
    primary key (id)
);

drop index participant_effects_sale_participations_idx;
