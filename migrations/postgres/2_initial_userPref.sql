-- +migrate Up
CREATE TABLE IF NOT EXISTS userpref (userId SERIAL PRIMARY KEY, settings TEXT NOT NULL);
-- +migrate Down
DROP TABLE userpref;