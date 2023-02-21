-- +migrate Up

create table matches (
  id bigint not null,
  operation_id bigint not null,
  offer_id bigint not null,
  base_amount numeric(20,0) not null,
  quote_amount numeric(20,0) not null,
  base_asset varchar(16) not null,
  quote_asset varchar(16) not null,
  price numeric(20,0) not null,
  created_at timestamp without time zone,

  primary key (base_asset, quote_asset, id)
);

-- +migrate Down

drop table if exists matches cascade;
