CREATE TABLE transaction_items
(
    "id"             SERIAL PRIMARY KEY,
    "transaction_id" TEXT,
    "product_id"     TEXT,
    "size_id"        TEXT,
    "flavor_id"      TEXT,
    "quantity"       INT,
    "subtotal"       DOUBLE PRECISION,
    "created_at"     TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at"     TIMESTAMP,
    "deleted_at"     TIMESTAMP,
    "created_by"     TEXT,
    "updated_by"     TEXT,
    "deleted_by"     TEXT,

    FOREIGN KEY ("transaction_id") REFERENCES transactions ("id"),
    FOREIGN KEY ("product_id") REFERENCES products ("id"),
    FOREIGN KEY ("size_id") REFERENCES sizes ("id"),
    FOREIGN KEY ("flavor_id") REFERENCES flavors ("id")
);