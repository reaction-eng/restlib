-- +migrate Up
CREATE TABLE IF NOT EXISTS resetrequests (id SERIAL PRIMARY KEY, userId int NOT NULL, email TEXT NOT NULL, token TEXT NOT NULL,issued DATE NOT NULL, type int NOT NULL);

-- +migrate Down
DROP TABLE resetrequests;