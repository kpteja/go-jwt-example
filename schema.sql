-- Create database "todos".
DROP DATABASE IF EXISTS todos;
CREATE DATABASE todos;
\c todos;

-- Create table "users" to store user details.
CREATE TABLE IF NOT EXISTS users (
    id SERIAL NOT NULL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    createdon TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create table "todos" to store todos created by users.
CREATE TABLE IF NOT EXISTS todos (
    id SERIAL NOT NULL PRIMARY KEY,
    todo TEXT NOT NULL,
    userid BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    createdon TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create table "refreshtokens" to store refresh tokens.
-- Refersh tokens are used to renew JWTs.
CREATE TABLE IF NOT EXISTS refreshtokens (
    id SERIAL NOT NULL PRIMARY KEY,
    refreshtoken TEXT NOT NULL UNIQUE,
    userid BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    expiry TIMESTAMP NOT NULL,
    createdon TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);