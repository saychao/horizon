-- +migrate Up


create table polls (
	id  bigint not null,
  permission_type int not null,
  number_of_choices int not null,
  data jsonb not null,
  start_time   TIMESTAMP without time zone NOT NULL,
  end_time     TIMESTAMP without time zone NOT NULL,
  owner_id text not null,
  result_provider_id text not null,
  vote_confirmation_required boolean not null,
  state int not null,
  details jsonb not null,
  primary key (id)
	);

create table votes (
  id bigint not null,
  poll_id bigint not null,
  voter_id text not null,
  data jsonb not null,
  primary key (id)
  );

create index by_poll_id ON votes USING btree (poll_id);

-- +migrate Down

drop table if exists polls cascade;
drop table if exists votes cascade;
