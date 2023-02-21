-- +migrate Up

ALTER TABLE reviewable_request ADD all_tasks        INT   NOT NULL DEFAULT 0;
ALTER TABLE reviewable_request ADD pending_tasks    INT   NOT NULL DEFAULT 0;
ALTER TABLE reviewable_request ADD external_details TEXT  NOT NULL DEFAULT '{}';

-- +migrate Down

ALTER TABLE reviewable_request DROP COLUMN all_tasks;
ALTER TABLE reviewable_request DROP COLUMN pending_tasks;
ALTER TABLE reviewable_request DROP COLUMN external_details;
