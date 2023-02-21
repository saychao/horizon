-- +migrate Up

alter table accounts add column kyc_recovery_status integer not null default 0;

-- +migrate Down

alter table accounts drop column if exists kyc_recovery_status;
