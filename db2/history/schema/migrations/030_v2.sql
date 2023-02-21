-- +migrate Up


create table ledgers (
    id bigint not null primary key,
    sequence bigint not null,
    hash character(64) not null,
    previous_hash character(64) not null,
    tx_count int not null,
    data text not null,
    closed_at timestamp without time zone
);

create table transactions (
    id bigint not null, --- consists ledger sequence as 32 significant bits. We should not worry to much about duplication
    --  of ids as such issue will be handled during ledger inserts.
    hash character(64) NOT NULL primary key,
    ledger_sequence integer NOT NULL,
    ledger_close_time timestamp without time zone,
    application_order integer NOT NULL,
    account character varying(64) NOT NULL,
    operation_count integer NOT NULL,
    envelope text NOT NULL,
    result text NOT NULL,
    meta text NOT NULL,
    valid_after timestamp without time zone not null,
    valid_before timestamp without time zone not null
);

create table ledger_changes (
  tx_id         bigint    not null,
  op_id         bigint    not null,
  order_number  int       not null,
  effect        int       not null,
  entry_type    int       not null,
  payload       text      not null
);

create index on ledger_changes (tx_id, entry_type);

create table accounts (
  id bigint not null primary key,
  address character varying(64) not null
);

create unique index on accounts (address);

create table balances (
  id bigint not null primary key,
  account_id bigint not null,
  address character varying(64) not null,
  asset_code character varying(64) not null
);

create unique index on balances (address);

create table operations (
  id bigint not null primary key,
  tx_id bigint not null,
  type int not null,
  details jsonb not null,
  ledger_close_time timestamp without time zone not null,
  source character varying(64) not null
);

create table participant_effects (
  id bigint not null primary key,
  account_id bigint not null,
  balance_id bigint ,
  asset_code character varying(64),
  effect jsonb,
  operation_id bigint not null
);

create index on participant_effects (balance_id, id) where balance_id is not null;

create table reviewable_requests (
  id bigint not null primary key,
  requestor character varying(64) not null,
  reviewer character varying(64) not null,
  reference character varying(64),
  reject_reason TEXT not null,
  request_type int not null,
  request_state int not null,
  hash character varying(64) not null,
  created_at timestamp without time zone not null,
  updated_at timestamp without time zone not null,
  details jsonb not null,
  all_tasks int not null,
  pending_tasks int not null,
  external_details jsonb not null
);

create table sales (
  id bigint not null primary key,
  soft_cap numeric(23, 0) not null,
  hard_cap numeric(23, 0) not null,
  base_current_cap numeric(23, 0) not null,
  base_hard_cap numeric(23, 0) not null,
  sale_type int not null,
  owner_address character varying(64) not null,
  base_asset character varying(64) not null,
  default_quote_asset character varying(64) not null,
  start_time timestamp without time zone not null,
  end_time timestamp without time zone not null,
  details jsonb not null,
  state int not null,
  quote_assets jsonb not null
);


-- +migrate Down

drop table ledgers cascade;
drop table transactions cascade;
drop table ledger_changes cascade;
drop table accounts cascade;
drop table balances cascade;
drop table operations cascade;
drop table participant_effects cascade;
drop table reviewable_requests cascade;
drop table sales cascade;
