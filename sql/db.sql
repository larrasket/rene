CREATE TABLE IF NOT EXISTS metadata
(
    id                INTEGER PRIMARY KEY,
    key               TEXT NOT NULL UNIQUE,
    value             TEXT NOT NULL,
    creation_date     TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
    modification_date TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime'))
);

CREATE TABLE IF NOT EXISTS accounts
(
    id             INTEGER PRIMARY KEY,
    username       TEXT NOT NULL UNIQUE,
    account_token  TEXT NOT NULL,
    account_secret TEXT NOT NULL
);

CREATE TRIGGER IF NOT EXISTS check_username_validity
    BEFORE INSERT
    ON accounts
BEGIN
    SELECT CASE
               WHEN new.username NOT LIKE '^[a-zA-Z0-9]+$' THEN
                   RAISE(ABORT, 'Username can only contain letters and numbers')
               END;
END;
