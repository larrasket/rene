CREATE TABLE IF NOT EXISTS metadata
(
    id    INTEGER PRIMARY KEY,
    key   TEXT NOT NULL UNIQUE,
    value TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS accounts
(
    id             INTEGER PRIMARY KEY,
    username       TEXT NOT NULL UNIQUE,
    account_token  TEXT NOT NULL,
    account_secret TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS queue
(
    id      INTEGER PRIMARY KEY,
    account TEXT,
    text    TEXT,
    CONSTRAINT fk_accounts FOREIGN KEY (account) REFERENCES accounts (username)
);
