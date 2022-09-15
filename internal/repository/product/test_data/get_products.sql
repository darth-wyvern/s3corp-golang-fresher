INSERT INTO users (id,name, email, phone, password)
    VALUES (1,'admin', 'admin@example.com', '0987654321', '123456789'),
    (2,'admin 2', 'admin2@example.com', '0987654321', '123456789');

INSERT INTO products (id,title,description, price, quantity, user_id)VALUES 
    (1,'Book : "Tam Cam"','The folk story of vietnamese, about a girl who was Tam', 100000, 1, 1), 
    (2,'Thien Long ballpoint pen','Nice pen from Thien Long company', 10000, 2, 1),
    (3,'Thien long pencil','Nice pen from Thien Long company', 5000, 3, 2),
    (4,'Book : "Dac nhan tam"','The favious book of the world', 100000, 4, 1),
    (5,'Note book','It is small and pretty', 20000, 5, 1),
    (6,'Sneaker','Beutyfull shoe', 5000000, 10, 1),
    (7,'Gucci backpack','The best of the backpacks', 6000000, 6, 1),
    (8,'Sneaker backpack','The best of the backpacks', 6500000, 7, 1),
    (9,'Gucci TShirt','Fashion TShirt', 1000000, 10, 2),
    (10,'Book : "Cha giau cha ngheo"','Nice to read every weekend', 200000, 8, 2),
    (11,'Gucci trousers','To be a gentleman', 10000000, 9, 2);
