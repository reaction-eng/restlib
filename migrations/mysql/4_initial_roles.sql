-- +migrate Up
CREATE TABLE IF NOT EXISTS roles (id int NOT NULL AUTO_INCREMENT, userId int, orgId int, roleId int, PRIMARY KEY (id) )

-- +migrate Down
DROP TABLE roles;
