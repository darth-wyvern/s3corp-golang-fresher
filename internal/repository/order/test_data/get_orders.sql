INSERT INTO "users" ("id", "name", "email", "password", "phone", "role", "is_active")
VALUES (10, 'test1', 'test1@example.com', 'test', 'test', 'ADMIN', true),
 (11, 'test2', 'test2@example.com', 'abcdef', '0987654321', 'ADMIN', true);

INSERT INTO "orders" ("id", "order_number", "user_id", "status")
VALUES (1, 'ORDER_NUMBER_1', 10, 'NEW'),
(2, 'ORDER_NUMBER_2', 11, 'NEW'),
(5, 'ORDER_NUMBER_5', 10, 'NEW'),
(4, 'ORDER_NUMBER_3', 11, 'NEW'),
(3, 'ORDER_NUMBER_4', 11, 'NEW');

INSERT INTO "products" ("id", "title", "description", "price", "quantity", "is_active", "user_id")
VALUES (10, 'Product 10', 'Product 10', 1000, 10, true, 10),
       (11, 'Product 11', 'Product 11', 2000, 11, true, 10);

INSERT INTO "order_items" ("id", "order_id", "product_id", "product_price", "product_name", "quantity", "discount", "note")
VALUES (10, 1, 10, 1000, 'Product 10', 20, 0, ''),
(11, 1, 11, 1000, 'Product 11', 30, 0, ''),
(12, 2, 11, 1000, 'Product 11', 40, 0, ''),
(13, 2, 10, 1000, 'Product 10', 50, 0, '');
