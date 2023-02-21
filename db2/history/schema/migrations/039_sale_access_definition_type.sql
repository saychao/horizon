-- +migrate Up

alter table sales add column access_definition_type int not null default 0;

-- +migrate Down

alter table sales drop column if exists access_definition_type;
