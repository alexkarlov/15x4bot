CREATE USER bot;
CREATE DATABASE bot;
GRANT ALL PRIVILEGES ON DATABASE bot TO bot;
\c bot bot;
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    starttime timestamp without time zone,
    endtime timestamp without time zone,
    name character varying(100),
    description character varying(1000),
    place integer
);

CREATE TABLE places (
    id SERIAL PRIMARY KEY,
    name character varying(100),    
    address character varying(100),
    description character varying(1000),
    map_url character varying(200)
);

CREATE TABLE actions (
    id SERIAL PRIMARY KEY,
    command character varying(30),
    msg character varying(2000)
);

CREATE TABLE repetitions (
    id SERIAL PRIMARY KEY,
    time timestamp,
    place int
);
