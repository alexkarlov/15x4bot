CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    starttime TIMESTAMP,
    endtime TIMESTAMP,
    description VARCHAR(1000),
    place INT
);

CREATE TABLE IF NOT EXISTS places (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),    
    address VARCHAR(100),
    description VARCHAR(1000),
    map_url VARCHAR(200)
);

CREATE TABLE IF NOT EXISTS actions (
    id SERIAL PRIMARY KEY,
    command VARCHAR(30),
    msg VARCHAR(2000)
);

CREATE TABLE IF NOT EXISTS repetitions (
    id SERIAL PRIMARY KEY,
    time TIMESTAMP,
    place INT
);
