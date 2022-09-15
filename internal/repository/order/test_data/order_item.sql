INSERT INTO "users" ("id", "name", "email", "password", "phone", "role", "is_active")
VALUES (10, 'test1', 'test1@example.com', 'test', 'test', 'ADMIN', true);

INSERT INTO "orders" ("id", "order_number", "order_date", "user_id", "status")
VALUES (10, 'AAA', '2022-08-04 02:00:00', 10, 'NEW');

INSERT INTO "products" ("id", "title", "description", "price", "quantity", "is_active", "user_id")
VALUES (10, 'product 10', 'Product 10', 1000, 10, true, 10),
       (11, 'product 11', 'Product 11', 2000, 11, true, 10);

INSERT INTO "order_items"
("id", "order_id", "product_id", "product_price", "product_name", "quantity", "discount", "note")
VALUES (10, 10, 10, 1000, 'product 10', 10, 0, '');
