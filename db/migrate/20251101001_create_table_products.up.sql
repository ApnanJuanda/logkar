CREATE TABLE products
(
    "id"         TEXT PRIMARY KEY,
    "name"       TEXT,
    "seller_id"  TEXT,
    "type_id"    TEXT,
    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP,
    "deleted_at" TIMESTAMP,
    "created_by" TEXT,
    "updated_by" TEXT,
    "deleted_by" TEXT,

    FOREIGN KEY ("seller_id") REFERENCES sellers ("id"),
    FOREIGN KEY ("type_id") REFERENCES product_types("id")
);