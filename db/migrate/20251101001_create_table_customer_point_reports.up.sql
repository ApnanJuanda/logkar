CREATE TABLE customer_point_reports (
    "id"                       SERIAL PRIMARY KEY,
    "customer_id"              TEXT,
    "transaction_id"           TEXT,
    "customer_point_redeem_id" INT,
    "status"                   TEXT,
    "balance"                  INT,
    "point_in"                 INT,
    "point_out"                INT,
    "created_at"               TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at"               TIMESTAMP,
    "deleted_at"               TIMESTAMP,
    "created_by"               TEXT,
    "updated_by"               TEXT,
    "deleted_by"               TEXT,

    FOREIGN KEY ("customer_id") REFERENCES customers ("id")
);
