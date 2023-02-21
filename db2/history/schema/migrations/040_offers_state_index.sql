-- +migrate Up

create index offers_state_indx on history_operations (type, (details->>'offer_id')) where type = 16;

-- +migrate Down

drop index offers_state_indx;