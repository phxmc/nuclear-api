BEGIN;

CREATE TABLE passes (
    id CHAR(42) PRIMARY KEY,
    account_id CHAR(42) REFERENCES accounts(id),
    "from" TIMESTAMP WITH TIME ZONE,
    "to" TIMESTAMP WITH TIME ZONE,
    active BOOLEAN,
    created_at TIMESTAMP WITH TIME ZONE
);

COMMIT;
