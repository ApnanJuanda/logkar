CREATE TABLE sellers (
    "id"                 TEXT PRIMARY KEY,
    "name"               TEXT      NOT NULL,
    "email"              TEXT      NOT NULL,
    "encrypted_password" TEXT      NOT NULL,
    "phone"              TEXT,
    "created_at"         TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at"         TIMESTAMP,
    "deleted_at"         TIMESTAMP,
    "created_by"         TEXT,
    "updated_by"         TEXT,
    "deleted_by"         TEXT
);