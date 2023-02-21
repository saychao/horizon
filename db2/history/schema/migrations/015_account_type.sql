-- +migrate Up

alter table history_accounts add column account_type int not null default 0;

-- +migrate Down

alter table history_accounts drop column account_type;