-- +migrate Up
CREATE TABLE IF NOT EXISTS userpref (userId int NOT NULL, settings TEXT NOT NULL, PRIMARY KEY (userId) );


-- +migrate Down
DROP TABLE userpref;