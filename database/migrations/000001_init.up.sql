CREATE TABLE IF NOT EXISTS dummy
(
    id         BYTEA PRIMARY KEY        NOT NULL,
    title      VARCHAR,
    status     INTEGER                  NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);