-- +migrate Up

CREATE INDEX hbu_balance_id_index ON history_balance_updates USING HASH (balance_id);

-- +migrate Down

DROP INDEX hbu_balance_id_index;
