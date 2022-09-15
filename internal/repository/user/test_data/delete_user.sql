INSERT INTO "users" ("id", "name", "email", "password", "phone", "role", "is_active")
VALUES (2, 'test1', 'test77@example.com', 'test', 'test', 'ADMIN', true),
       (1, 'test1', 'test1@example.com', 'test', 'test', 'ADMIN', true);

INSERT INTO "products" (id, title, description, price, quantity, is_active, user_id)
VALUES (1, 'test1', 'test1', 1, 1, true, 2);
