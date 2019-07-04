CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    starttime timestamp without time zone,
    endtime timestamp without time zone,
    name VARCHAR(100),
    description VARCHAR(1000),
    place integer
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
    time timestamp,
    place int
);
