-- +goose Up
CREATE TABLE feeds(
     id int PRIMARY KEY,
     created_at TIMESTAMP NOT NULL,
     updated_at TIMESTAMP NOT NULL,
     name TEXT NOT NULL,
     url TEXT  UNIQUE NOT NULL,
     user_id int NOT NULL 
     REFERENCES users
     ON DELETE CASCADE

);
-- +goose Down
DROP TABLE feeds;