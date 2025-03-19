BEGIN;

ALTER TABLE accounts
DROP COLUMN telegram_id;

COMMIT;
