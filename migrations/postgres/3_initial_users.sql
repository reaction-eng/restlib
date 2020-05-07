-- +migrate Up
CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, email TEXT NOT NULL, password TEXT NOT NULL, activation Date)

-- +migrate Down
DROP TABLE users;