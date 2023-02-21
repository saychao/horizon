-- +migrate Up


create table asset_pairs (
		base character varying(64) not null,
		quote character varying(64) not null,
		current_price  bigint       not null,
		ledger_close_time timestamp without time zone not null
	);

-- +migrate Down

drop table asset_pairs cascade;
