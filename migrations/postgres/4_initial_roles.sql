-- +migrate Up
CREATE TABLE IF NOT EXISTS roles (id SERIAL PRIMARY KEY, userId int NOT NULL, orgId int NOT NULL, roleId int NOT NULL);

-- +migrate Down
DROP TABLE roles;