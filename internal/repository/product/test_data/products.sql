INSERT INTO users ("id", "name", "email", "phone", password)
VALUES (1, 'admin', 'admin@example.com', '0987654321', '123456789');

INSERT INTO products ("id", "title", "price", "quantity", "user_id", "is_active")
VALUES (1, 'AAA', 20000, 10, 1, true),
       (2, 'BBB', 15000, 20, 1, false),
       (3, 'CCC', 23000, 62, 1, true),
       (4, 'DDD', 28000, 13, 1, true);