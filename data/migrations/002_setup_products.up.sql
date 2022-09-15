BEGIN;

CREATE TABLE IF NOT EXISTS "products" (
    "id" SERIAL PRIMARY KEY,
    "title" TEXT NOT NULL,
    "description" TEXT NOT NULL DEFAULT '',
    "price" FLOAT NOT NULL,
    "quantity" INT NOT NULL DEFAULT 0,
    "is_active" BOOL NOT NULL DEFAULT TRUE,
    "user_id" INT NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY ("user_id") REFERENCES "users"("id")
);

COMMIT;
