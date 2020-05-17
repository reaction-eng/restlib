-- +migrate Up
CREATE TABLE IF NOT EXISTS users (id int NOT NULL AUTO_INCREMENT, email TEXT, password TEXT, activation Date, PRIMARY KEY (id) )
CREATE TABLE IF NOT EXISTS userOrganizations (userId int, orgId int, joinDate Date )

-- +migrate Down
DROP TABLE users;
DROP TABLE userpref;