CREATE TABLE transactions (
    "id"           TEXT PRIMARY KEY,
    "customer_id"  TEXT,
    "total_amount" DOUBLE PRECISION,
    "status"       INT,
    "note"         TEXT,
    "created_at"   TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at"   TIMESTAMP,
    "deleted_at"   TIMESTAMP,
    "created_by"   TEXT,
    "updated_by"   TEXT,
    "deleted_by"   TEXT,

    FOREIGN KEY ("customer_id") REFERENCES customers("id"),
    FOREIGN KEY ("status") REFERENCES transaction_statuses("status_code")
);
