-- +migrate Up

create table asset
(
    code                    varchar(16)                 not null,
    owner                   varchar(56)                 not null,
    preissued_asset_signer  varchar(56)                 not null,

    max_issuance_amount     numeric(20,0)               not null,
    issued                  numeric(20,0)               not null,
    available_for_issuance  numeric(20,0)               not null,
    pending_issuance        numeric(20,0)               not null,
    type                    bigint                      not null,
    policies                int                         not null,

    details                 jsonb                       not null,
    trailing_digits         int                         not null,
    state                   int                         not null,
    primary key (code)
);

create index asset_by_owner on asset using btree (owner);
create index asset_by_preissued_asset_signer on asset using btree (preissued_asset_signer);

-- +migrate Down

drop table if exists asset cascade;