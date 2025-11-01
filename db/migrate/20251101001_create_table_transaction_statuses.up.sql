CREATE TABLE transaction_statuses (
    "id"          TEXT PRIMARY KEY,
    "name"        TEXT,
    "label"       TEXT,
    "status_code" INT UNIQUE,
    "created_at"  TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at"  TIMESTAMP,
    "deleted_at"  TIMESTAMP,
    "created_by"  TEXT,
    "updated_by"  TEXT,
    "deleted_by"  TEXT
);
