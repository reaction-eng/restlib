-- +migrate Up
CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, email TEXT NOT NULL, password TEXT NOT NULL, activation Date)
CREATE TABLE IF NOT EXISTS userpref (userId int NOT NULL, orgId int NOT NULL, joinDate Date)

-- +migrate Down
DROP TABLE users;
DROP TABLE userpref;