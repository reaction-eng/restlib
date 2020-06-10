-- +migrate Up
CREATE TABLE IF NOT EXISTS resetrequests (id int NOT NULL AUTO_INCREMENT, userId int, email TEXT, token TEXT, issued DATE, type INT, PRIMARY KEY (id) );


-- +migrate Down
DROP TABLE resetrequests;