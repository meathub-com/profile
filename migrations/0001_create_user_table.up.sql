CREATE TABLE IF NOT EXISTS sellers
(
    id       SERIAL PRIMARY KEY,
    name     VARCHAR(100) NOT NULL,
    user_id  INT NOT NULL
)

CREATE TABLE IF NOT EXISTS addresses
(
    id          SERIAL PRIMARY KEY,
    seller_id   INT NOT NULL,
    street      VARCHAR(200) NOT NULL,
    city        VARCHAR(100) NOT NULL,
    state       VARCHAR(100) NOT NULL,
    postal_code VARCHAR(20) NOT NULL,
    country     VARCHAR(100) NOT NULL,
    FOREIGN KEY (seller_id) REFERENCES sellers (id)
    );
