-- Create tables orders and orders items and create unique index for each.
BEGIN;

CREATE TABLE IF NOT EXISTS "orders"
(
    "id" SERIAL PRIMARY KEY,
    "order_number" TEXT NOT NULL,
    "order_date" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "status" TEXT NOT NULL DEFAULT 'NEW', -- NEW, PENDING, SUCCESS, FAILED
    "note" TEXT NOT NULL DEFAULT '',
    "user_id" INT NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY ("user_id") REFERENCES "users"("id")
);

CREATE UNIQUE INDEX IF NOT EXISTS "order_number_on_orders" ON "orders"("order_number");

CREATE TABLE IF NOT EXISTS "order_items" (
    "id" SERIAL PRIMARY KEY,
    "order_id" INT NOT NULL,
    "product_id" INT NOT NULL,
    "product_price" FLOAT NOT NULL,
    "product_name" TEXT NOT NULL,
    "quantity" INT NOT NULL,
    "discount" FLOAT NOT NULL DEFAULT 0,
    "note" TEXT NOT NULL DEFAULT '',
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY ("order_id") REFERENCES "orders"("id"),
    FOREIGN KEY ("product_id") REFERENCES "products"("id")
);

CREATE UNIQUE INDEX IF NOT EXISTS "order_id_product_id_on_order_items" ON "order_items"("order_id", "product_id");

END;
