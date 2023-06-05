CREATE TABLE IF NOT EXISTS profiles
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255),
    user_id    VARCHAR(255)
);
CREATE TABLE IF NOT EXISTS address
(
    id          SERIAL PRIMARY KEY,
    profile_id  INTEGER REFERENCES profiles(id),
    street      VARCHAR(255),
    city        VARCHAR(255),
    state       VARCHAR(255),
    postal_code VARCHAR(255),
    country     VARCHAR(255)
);


INSERT INTO profiles (name, user_id)
VALUES ('John Doe', 'johndoe123'),
       ('Jane Smith', 'janesmith456'),
       ('Michael Johnson', 'mjohnson789');

INSERT INTO address (profile_id, street, city, state, postal_code, country)
VALUES (1, '123 Main Street', 'New York', 'NY', '10001', 'United States'),
       (2, '456 Elm Street', 'Los Angeles', 'CA', '90001', 'United States'),
       (3, '789 Oak Street', 'Chicago', 'IL', '60601', 'United States');