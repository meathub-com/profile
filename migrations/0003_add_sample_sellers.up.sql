-- insert data into sellers
INSERT INTO sellers (name, user_id)
VALUES
    ('Seller1', 1),
    ('Seller2', 2),
    ('Seller3', 3),
    ('Seller4', 4),
    ('Seller5', 5);

-- insert data into addresses
INSERT INTO addresses (seller_id, street, city, state, postal_code, country)
VALUES
    (1, '123 Main Street', 'San Francisco', 'CA', '94101', 'USA'),
    (2, '456 Maple Drive', 'Chicago', 'IL', '60601', 'USA'),
    (3, '789 Oak Street', 'New York', 'NY', '10001', 'USA'),
    (4, '321 Elm Street', 'Los Angeles', 'CA', '90001', 'USA'),
    (5, '654 Pine Avenue', 'Austin', 'TX', '73301', 'USA');
