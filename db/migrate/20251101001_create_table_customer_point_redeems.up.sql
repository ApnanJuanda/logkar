CREATE TABLE customer_point_redeems (
    "id"                 SERIAL PRIMARY KEY,
    "customer_id"        TEXT,
    "total_redeem_point" INT,
    "created_at"         TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at"         TIMESTAMP,
    "deleted_at"         TIMESTAMP,
    "created_by"         TEXT,
    "updated_by"         TEXT,
    "deleted_by"         TEXT,

    FOREIGN KEY ("customer_id") REFERENCES customers ("id")
);
