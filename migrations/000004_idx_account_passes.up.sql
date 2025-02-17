BEGIN;

CREATE UNIQUE INDEX idx_account_passes
ON account_passes (account_id)
WHERE is_active = TRUE;

COMMIT;
