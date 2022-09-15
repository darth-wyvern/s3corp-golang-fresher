INSERT INTO "users" ("id", "name", "email", "password", "phone", "role", "is_active")
VALUES (10, 'test1', 'test1@example.com', 'test', 'test', 'ADMIN', true);

INSERT INTO "orders" ("id", "order_number", "user_id", "status")
VALUES (10, 'AAA', 10, 'NEW'),
       (11, 'BBB', 10, 'NEW'),
       (12, 'CCC', 10, 'FAILED'),
       (13, 'DDD', 10, 'SUCCESS'),
       (14, 'EEE', 10, 'PENDING'),
       (15, 'FFF', 10, 'SUCCESS');
