-- +migrate Up


create table support_params
(
	ingest_version int not null check (ingest_version >= 0)
);

-- +migrate Down

drop table if exists support_params;
