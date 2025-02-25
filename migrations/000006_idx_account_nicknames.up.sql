BEGIN;

CREATE UNIQUE INDEX idx_unique_active_nickname
ON account_nicknames (nickname)
WHERE is_active = TRUE;

CREATE UNIQUE INDEX idx_account_nicknames_active
ON account_nicknames (account_id)
WHERE is_active = TRUE;

COMMIT;
