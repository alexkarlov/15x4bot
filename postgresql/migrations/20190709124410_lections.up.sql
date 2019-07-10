CREATE TABLE IF NOT EXISTS lections (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    description VARCHAR(1000),
    user_id INT NOT NULL
);
