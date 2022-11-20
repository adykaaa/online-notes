DROP DATABASE IF EXISTS users;
CREATE DATABASE users_db;

CREATE TABLE IF NOT EXISTS users(
  id serial PRIMARY KEY,
  email VARCHAR(50) NOT NULL UNIQUE,
  username VARCHAR(50) NOT NULL UNIQUE,
  password VARCHAR(50) NOT NULL
);
