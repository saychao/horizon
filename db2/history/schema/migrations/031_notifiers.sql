-- +migrate Up
--
-- Must stay as one line as migrator treats it as several separate statements
--
CREATE OR REPLACE FUNCTION new_history_ledger_seq_notify() RETURNS trigger AS $$ DECLARE payload varchar; mid uuid; BEGIN payload = CAST(NEW.sequence AS text); PERFORM pg_notify('new_history_ledger_seq', payload); RETURN NEW; END; $$ LANGUAGE plpgsql;

CREATE TRIGGER history_ledgers_insert AFTER INSERT ON history_ledgers FOR EACH ROW EXECUTE PROCEDURE new_history_ledger_seq_notify();


CREATE OR REPLACE FUNCTION new_ledgers_seq_notify() RETURNS trigger AS $$ DECLARE payload varchar; mid uuid; BEGIN payload = CAST(NEW.sequence AS text); PERFORM pg_notify('new_ledgers_seq', payload); RETURN NEW; END; $$ LANGUAGE plpgsql;

CREATE TRIGGER ledgers_insert AFTER INSERT ON ledgers FOR EACH ROW EXECUTE PROCEDURE new_ledgers_seq_notify();

-- +migrate Down

DROP TRIGGER IF EXISTS history_ledgers_insert ON history_ledgers;
DROP FUNCTION IF EXISTS new_history_ledger_seq_notify(varchar, uuid);
DROP TRIGGER IF EXISTS ledgers_insert ON ledgers;
DROP FUNCTION IF EXISTS new_ledgers_seq_notify(varchar, uuid);