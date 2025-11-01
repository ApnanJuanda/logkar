CREATE TABLE product_details (
    "id"         SERIAL PRIMARY KEY,
    "product_id" TEXT,
    "size_id"    TEXT,
    "flavor_id"  TEXT,
    "price"      DOUBLE PRECISION   DEFAULT 0,
    "stock"      INT,
    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP,
    "deleted_at" TIMESTAMP,
    "created_by" TEXT,
    "updated_by" TEXT,
    "deleted_by" TEXT,

    FOREIGN KEY ("product_id") REFERENCES products ("id"),
    FOREIGN KEY ("size_id") REFERENCES sizes ("id"),
    FOREIGN KEY ("flavor_id") REFERENCES flavors ("id")
);
