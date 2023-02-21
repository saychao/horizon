-- +migrate Up

create table history_balance_updates (
	id bigserial primary key,
	balance_id varchar(64) references history_balances (balance_id) not null,
	amount bigint not null,
	updated_at timestamp with time zone not null -- time zone is intentional, for easier JSON parsing
);

-- +migrate Down

drop table history_balance_updates;