CREATE TABLE IF NOT EXISTS sellers
(
    id      SERIAL PRIMARY KEY,
    name    VARCHAR(100) NOT NULL,
    user_id INT          NOT NULL
);

CREATE TABLE IF NOT EXISTS addresses
(
    id          SERIAL PRIMARY KEY,
    seller_id   INT          NOT NULL,
    street      VARCHAR(200) NOT NULL,
    city        VARCHAR(100) NOT NULL,
    state       VARCHAR(100) NOT NULL,
    postal_code VARCHAR(20)  NOT NULL,
    country     VARCHAR(100) NOT NULL,
    FOREIGN KEY (seller_id) REFERENCES sellers (id)
);
INSERT INTO sellers (name, user_id)
VALUES ('BlueSky Enterprises', 1),
       ('Elite Solutions Inc.', 2),
       ('Global Trade Co.', 3);

INSERT INTO addresses (seller_id, street, city, state, postal_code, country)
VALUES (1, '123 Main Street', 'City 1', 'State 1', '12345', 'Country 1'),
       (2, '456 Elm Street', 'City 2', 'State 2', '67890', 'Country 2'),
       (3, '789 Oak Street', 'City 3', 'State 3', '54321', 'Country 3');