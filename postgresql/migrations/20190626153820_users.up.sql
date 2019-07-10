CREATE TYPE userrole AS ENUM ('admin','lector', 'guest');
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100),
    role userrole DEFAULT 'guest',
    name VARCHAR(100),
    fb VARCHAR(100),
    vk VARCHAR(100),
    bdate TIMESTAMP,
    cdate TIMESTAMP NOT NULL DEFAULT NOW(),
    udate TIMESTAMP NOT NULL DEFAULT NOW()
);