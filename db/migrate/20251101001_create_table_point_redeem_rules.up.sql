CREATE TABLE point_redeem_rules (
    "id"             SERIAL PRIMARY KEY,
    "exchange_point" INT,
    "size_id"        TEXT,
    "created_at"     TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at"     TIMESTAMP,
    "deleted_at"     TIMESTAMP,
    "created_by"     TEXT,
    "updated_by"     TEXT,
    "deleted_by"     TEXT,

    FOREIGN KEY ("size_id") REFERENCES sizes ("id")
);
