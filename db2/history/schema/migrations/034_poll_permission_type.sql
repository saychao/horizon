-- +migrate Up
ALTER TABLE polls ALTER permission_type TYPE NUMERIC(10,0);

-- +migrate Down

-- we don't change type back to int because max numeric(10,0) greater than max int
-- 2019/04/09 18:42:16 pq: integer out of range handling 034_poll_permission_type.sql
